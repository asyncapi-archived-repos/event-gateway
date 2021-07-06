//nolint
package protocol

// The types on this file have been copied from https://github.com/Shopify/sarama and are used for decoding requests.
// As the decoder/encoder interfaces in Sarama project are not public, there is no way of reusing.
// The following issue has been opened https://github.com/Shopify/sarama/issues/1967.

import "github.com/klauspost/compress/zstd"

var (
	zstdDec, _ = zstd.NewReader(nil)
)

func zstdDecompress(dst, src []byte) ([]byte, error) {
	return zstdDec.DecodeAll(src, dst)
}
