package util

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// TODO - generalize further
func ParseMarkdownFirstSectionAsTviewText(content []byte) (s string, err error) { // no point in passing a Reader, because we need random access

	p := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithHeadingAttribute(),
		),
	)

	doc := p.Parser().Parse(text.NewReader(content))

	var bs []byte
	var attrs tcell.AttrMask
	err = ast.Walk(doc, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		switch node.(type) {
		case *ast.Document:
		case *ast.CodeSpan:
			if entering {
				attrs |= tcell.AttrItalic
				bs = appendAttributes(bs, attrs)
			} else {
				attrs = attrs &^ tcell.AttrItalic
				bs = appendAttributes(bs, attrs)
			}
		case *ast.Emphasis:
			if entering {
				attrs |= tcell.AttrBold
				bs = appendAttributes(bs, attrs)
			} else {
				attrs = attrs &^ tcell.AttrBold
				bs = appendAttributes(bs, attrs)
			}
		case *ast.Heading:
			if heading, ok := node.(*ast.Heading); ok && heading.Level != 1 {
				return ast.WalkStop, nil
			}
			if entering {
				attrs |= tcell.AttrUnderline | tcell.AttrBold | tcell.AttrItalic
				bs = appendAttributes(bs, attrs)
			} else {
				attrs = attrs &^ tcell.AttrUnderline &^ tcell.AttrBold &^ tcell.AttrItalic
				bs = appendAttributes(bs, attrs)
				bs = append(bs, '\n')
			}
		case *ast.Link:
			if entering {
				attrs |= tcell.AttrUnderline
				bs = appendAttributes(bs, attrs)
			} else {
				attrs = attrs &^ tcell.AttrUnderline
				bs = appendAttributes(bs, attrs)
			}
		case *ast.List:
			if entering {
				bs = append(bs, '\n')
			}
		case *ast.ListItem:
			if entering {
				bs = append(bs, ' ')
				attrs |= tcell.AttrBold
				bs = appendAttributes(bs, attrs)
				bs = append(bs, '*')
				attrs = attrs &^ tcell.AttrBold
				bs = appendAttributes(bs, attrs)
				bs = append(bs, ' ')
			} else {
				bs = append(bs, '\n')
			}
		case *ast.Paragraph:
			bs = append(bs, '\n')
		case *ast.TextBlock:
		case *ast.Text:
			if entering {
				bs = append(bs, node.Text(content)...)
			} else {
				if node.NextSibling() != nil && bs[len(bs)-1] != ' ' {
					bs = append(bs, ' ')
				}
			}
		default:
			panic(fmt.Sprintf("unsupported node type: %T", node))
		}
		return ast.WalkContinue, nil
	})
	if err != nil {
		return
	}
	s = string(bs)
	return
}

func appendAttributes(bs []byte, attrs tcell.AttrMask) (result []byte) {
	result = append(bs, '[', ':', ':')
	if attrs&tcell.AttrBold != 0 {
		result = append(result, 'b')
	}
	if attrs&tcell.AttrItalic != 0 {
		result = append(result, 'i')
	}
	if attrs&tcell.AttrUnderline != 0 {
		result = append(result, 'u')
	}
	if attrs == 0 {
		result = append(result, '-')
	}
	if attrs&^(tcell.AttrBold|tcell.AttrItalic|tcell.AttrUnderline) > 0 {
		panic(fmt.Sprintf("unsupported attribute: %d", attrs))
	}
	result = append(result, ']')
	return
}
