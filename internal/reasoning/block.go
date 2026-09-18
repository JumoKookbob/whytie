package reasoning

import (
	"github.com/JumoKookbob/whytie/internal/scanner"
	"github.com/JumoKookbob/whytie/internal/syntax"
)

type Item struct {
	Comment scanner.SourceComment
	Reasons []scanner.SourceComment
}

type Block struct {
	Items []scanner.SourceComment
}

func Group(comments []scanner.SourceComment) []Block {
	if len(comments) == 0 {
		return nil
	}

	var blocks []Block

	current := Block{
		Items: []scanner.SourceComment{comments[0]},
	}

	for i := 1; i < len(comments); i++ {
		previous := comments[i-1]
		comment := comments[i]

		sameFile := comment.File == previous.File
		adjacentLine := comment.Line == previous.Line+1

		if sameFile && adjacentLine {
			current.Items = append(current.Items, comment)
			continue
		}

		blocks = append(blocks, current)

		current = Block{
			Items: []scanner.SourceComment{comment},
		}
	}

	blocks = append(blocks, current)

	return blocks
}

func AttachReasons(block Block) []Item {
	var items []Item

	for _, comment := range block.Items {
		if comment.Kind == syntax.Reason {
			if len(items) == 0 {
				items = append(items, Item{
					Comment: comment,
				})
				continue
			}

			last := len(items) - 1

			if items[last].Comment.Kind == syntax.Reason {
				items = append(items, Item{
					Comment: comment,
				})
				continue
			}

			items[last].Reasons = append(
				items[last].Reasons,
				comment,
			)
			continue
		}

		items = append(items, Item{
			Comment: comment,
		})
	}

	return items
}
