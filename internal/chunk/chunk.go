package chunk

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

type Meta struct {
	IsMeta  bool   `json:"meta"`
	Version int    `json:"v"`
	DocID   string `json:"doc"`
	Hash    string `json:"hash"`
	Chunks  int    `json:"chunks"`
	Model   string `json:"model"`
	Dim     int    `json:"dim"`
	Mtime   int64  `json:"mtime"`
}

type Chunk struct {
	Idx    int       `json:"idx"`
	Title  string    `json:"title"`
	Offset int       `json:"offset"`
	Length int       `json:"length"`
	Vec    []float32 `json:"vec,omitempty"`
	Text   string    `json:"-"`
}

type ChunkType string

const (
	TypeProse ChunkType = "prose"
	TypeCode  ChunkType = "code"
	TypeTable ChunkType = "table"
	TypeList  ChunkType = "list"
)

type ChunkWithMeta struct {
	Meta   Meta
	Chunks []Chunk
}

func Split(content, docID string, maxLen, overlap int) ChunkWithMeta {
	h := sha256.Sum256([]byte(content))
	hash := hex.EncodeToString(h[:])[:16]

	chunks := splitIntoChunks(content, maxLen, overlap)

	return ChunkWithMeta{
		Meta: Meta{
			IsMeta:  true,
			Version: 1,
			DocID:   docID,
			Hash:    hash,
			Chunks:  len(chunks),
		},
		Chunks: chunks,
	}
}

func splitIntoChunks(content string, maxLen, overlap int) []Chunk {
	if len(content) <= maxLen {
		return []Chunk{{
			Idx:    0,
			Title:  "Overview",
			Offset: 0,
			Length: len(content),
		}}
	}

	var chunks []Chunk
	idx := 0

	separator := "\n## "
	sections := strings.Split(content, separator)
	pos := 0

	for i, sec := range sections {
		secStart := pos
		secLen := len(sec)
		pos += secLen + len(separator)

		if i == 0 && !strings.HasPrefix(content, "## ") {
			if secLen > 100 {
				chunks = append(chunks, Chunk{
					Idx:    idx,
					Title:  "Overview",
					Offset: secStart,
					Length: secLen,
				})
				idx++
			}
			continue
		}

		newlineIdx := strings.Index(sec, "\n")
		var title, body string
		if newlineIdx == -1 {
			title = sec
			body = ""
		} else {
			title = sec[:newlineIdx]
			body = sec[newlineIdx+1:]
		}

		bodyStart := secStart + len(title) + 1

		if len(body) > maxLen {
			subChunks := splitByParagraphs(title, body, maxLen, overlap, &idx)
			for j := range subChunks {
				subChunks[j].Offset = bodyStart + subChunks[j].Offset
			}
			chunks = append(chunks, subChunks...)
		} else if len(body) > 0 {
			chunks = append(chunks, Chunk{
				Idx:    idx,
				Title:  title,
				Offset: bodyStart,
				Length: len(body),
			})
			idx++
		}
	}

	return chunks
}

func splitByParagraphs(title, body string, maxLen, overlap int, idx *int) []Chunk {
	var chunks []Chunk
	localIdx := 0

	pos := 0
	chunkStart := 0
	chunkLen := 0

	for {
		nextPara := strings.Index(body[pos:], "\n\n")
		var paraEnd int
		if nextPara == -1 {
			paraEnd = len(body)
		} else {
			paraEnd = pos + nextPara
		}

		originalParaLen := paraEnd - pos
		para := strings.TrimSpace(body[pos:paraEnd])
		paraLen := len(para)

		if paraLen > 0 {
			if chunkLen > 0 && chunkLen+2+paraLen > maxLen {
				chunks = append(chunks, Chunk{
					Idx:    *idx,
					Title:  formatTitle(title, localIdx),
					Offset: chunkStart,
					Length: chunkLen,
				})
				*idx++
				localIdx++

				overlapStart := chunkStart + chunkLen - overlap
				if overlapStart < chunkStart {
					overlapStart = chunkStart
				}
				chunkStart = overlapStart
				chunkLen = paraEnd - overlapStart
			} else {
				if chunkLen == 0 {
					chunkStart = pos
				}
				if chunkLen > 0 {
					chunkLen += 2
				}
				chunkLen += originalParaLen
			}
		}

		if nextPara == -1 {
			break
		}
		pos = paraEnd + 2
	}

	if chunkLen > 0 {
		chunks = append(chunks, Chunk{
			Idx:    *idx,
			Title:  formatTitle(title, localIdx+1),
			Offset: chunkStart,
			Length: chunkLen,
		})
		*idx++
	}

	return chunks
}

func formatTitle(title string, part int) string {
	if part == 0 {
		return title
	}
	return title + " (" + string(rune(part+1)) + ")"
}

func getOverlap(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[len(s)-n:]
}

func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func ComputeHash(content string) string {
	h := sha256.Sum256([]byte(content))
	return hex.EncodeToString(h[:])[:16]
}
