package rs

import (
	"dsxy/objectstream"
	"fmt"
	"github.com/klauspost/reedsolomon"
	"io"
)

type RSGetStream struct {
	*decoder
}

type decoder struct {
	readers   []io.Reader
	writers   []io.Writer
	enc       reedsolomon.Encoder
	size      int64
	cache     []byte
	cacheSize int
	total     int64
}

func NewDecoder(readers []io.Reader, writers []io.Writer, size int64) *decoder {
	enc, _ := reedsolomon.New(DATA_SHARDS, PARITY_SHARDS)
	return &decoder{readers, writers, enc, size, nil, 0, 0}
}
func NewRSGetStream(localInfo map[int]string, dataServers []string, hash string, size int64) (*RSGetStream, error) {
	if len(localInfo)+len(dataServers) != ALL_SHARDS {
		return nil, fmt.Errorf("dataServers number mismatch")
	}
	readers := make([]io.Reader, ALL_SHARDS)
	for i := 0; i < ALL_SHARDS; i++ {
		server := localInfo[i]
		if server == "" {
			localInfo[i] = dataServers[0]
			dataServers = dataServers[1:]
			continue
		}
		if reader, e := objectstream.NewGetStream(server,
			fmt.Sprintf("%s.%d", hash, i)); e == nil {
			readers[i] = reader
		}
	}
	//Recover
	writers := make([]io.Writer, ALL_SHARDS)
	perShard := (size + DATA_SHARDS - 1) / DATA_SHARDS
	var e error
	for i := range readers {
		if readers[i] == nil {
			writers[i], e = objectstream.NewTempPutStream(localInfo[i],
				fmt.Sprintf("%s.%d", hash, i), perShard)
			if e != nil {
				return nil, e
			}
		}
	}
	//Init
	dec := NewDecoder(readers, writers, size)
	return &RSGetStream{dec}, nil

}
func (d *decoder) getData() error {
	if d.total == d.size {
		return io.EOF
	}
	shards := make([][]byte, ALL_SHARDS)
	repairIds := make([]int, 0)
	for i := range shards {
		if d.readers[i] == nil {
			repairIds = append(repairIds, i)
		} else {
			shards[i] = make([]byte, BLOCK_PER_SHARD)
			n, e := io.ReadFull(d.readers[i], shards[i])
			if e != nil && e != io.EOF && e != io.ErrUnexpectedEOF {
				shards[i] = nil
			} else if n != BLOCK_PER_SHARD {
				shards[i] = shards[i][:n]
			}
		}
	}
	e := d.enc.Reconstruct(shards)
	if e != nil {
		return e
	}
	//Recover(Write reconstructed shards to damaged writers)
	for _, idx := range repairIds {
		d.writers[idx].Write(shards[idx])
	}

	for i := 0; i < DATA_SHARDS; i++ {
		shardSize := int64(len(shards))
		if d.total+shardSize > d.size {
			reduce := d.total + shardSize - d.size
			shardSize -= reduce
		}
		d.cache = append(d.cache, shards[i][:shardSize]...)
		d.cacheSize += int(shardSize)
		d.total += shardSize
	}
	return nil
}
func (d *decoder) Read(p []byte) (n int, err error) {
	if d.cacheSize == 0 {
		if e := d.getData(); e != nil {
			return 0, e
		}
	}
	length := len(p)
	if d.cacheSize < length {
		length = d.cacheSize
	}
	d.cacheSize -= length
	copy(p, d.cache[:length])
	d.cache = d.cache[length:]
	return length, nil
}

func (s *RSGetStream) Close() {
	for i := range s.writers {
		if s.writers[i] != nil {
			s.writers[i].(*objectstream.TempPutStream).Commit(true)
		}
	}
}
