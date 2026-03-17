package chunk

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"unicode/utf8"
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
	if utf8.RuneCountInString(content) <= maxLen {
		return []Chunk{{
			Idx:    0,
			Title:  "Overview",
			Offset: 0,
			Length: utf8.RuneCountInString(content),
		}}
	}

	var chunks []Chunk
	offset := 0
	idx := 0

	sections := strings.Split(content, "\n## ")

	for i, sec := range sections {
		if i == 0 && !strings.HasPrefix(sec, "## ") {
			if len(sec) > 100 {
				chunks = append(chunks, Chunk{
					Idx:    idx,
					Title:  "Overview",
					Offset: offset,
					Length: utf8.RuneCountInString(sec),
				})
				idx++
			}
			offset += utf8.RuneCountInString(sec) + 3
			continue
		}

		lines := strings.SplitN(sec, "\n", 2)
		title := strings.TrimPrefix(lines[0], "## ")

		var body string
		if len(lines) > 1 {
			body = strings.TrimSpace(lines[1])
		}

		if utf8.RuneCountInString(body) > maxLen {
			subChunks := splitByParagraphs(title, body, maxLen, overlap, &idx)
			for _, sc := range subChunks {
				sc.Offset = offset
				offset += sc.Length
				chunks = append(chunks, sc)
			}
		} else {
			chunks = append(chunks, Chunk{
				Idx:    idx,
				Title:  title,
				Offset: offset,
				Length: utf8.RuneCountInString(body),
			})
			idx++
			offset += utf8.RuneCountInString(body) + 3
		}
	}

	return chunks
}

func splitByParagraphs(title, body string, maxLen, overlap int, idx *int) []Chunk {
	paras := strings.Split(body, "\n\n")
	var chunks []Chunk
	var current strings.Builder
	localIdx := 0

	for _, para := range paras {
		para = strings.TrimSpace(para)
		if para == "" {
			continue
		}

		if current.Len()+len(para) > maxLen && current.Len() > 0 {
			chunks = append(chunks, Chunk{
				Idx:    *idx,
				Title:  formatTitle(title, localIdx),
				Length: utf8.RuneCountInString(current.String()),
			})
			*idx++
			localIdx++

			overlapRunes := getOverlap(current.String(), overlap)
			current.Reset()
			current.WriteString(overlapRunes)
			current.WriteString("\n\n")
		}

		current.WriteString(para)
		current.WriteString("\n\n")
	}

	if current.Len() > 0 {
		chunks = append(chunks, Chunk{
			Idx:    *idx,
			Title:  formatTitle(title, localIdx+1),
			Length: utf8.RuneCountInString(strings.TrimSpace(current.String())),
		})
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
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[len(runes)-n:])
}

func Truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}

func ComputeHash(content string) string {
	h := sha256.Sum256([]byte(content))
	return hex.EncodeToString(h[:])[:16]
}
