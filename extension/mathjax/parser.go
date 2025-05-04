package mathjax

import (
	"bytes"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

type mathjaxInlineParser struct{}

func NewMathjaxInlineParser() parser.InlineParser {
	return &mathjaxInlineParser{}
}

func (w *mathjaxInlineParser) Trigger() []byte {
	return []byte{'$'}
}

func (p *mathjaxInlineParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, startSegment := block.PeekLine()
	opener := 0
	for ; opener < len(line) && line[opener] == '$'; opener++ {
	}
	block.Advance(opener)
	l, pos := block.Position()
	node := NewInlineMath()
	for {
		line, segment := block.PeekLine()
		if line == nil {
			block.SetPosition(l, pos)
			return ast.NewTextSegment(startSegment.WithStop(startSegment.Start + opener))
		}
		for i := 0; i < len(line); i++ {
			c := line[i]
			if c == '$' {
				oldi := i
				for ; i < len(line) && line[i] == '$'; i++ {
				}
				closure := i - oldi
				if closure == opener && (i >= len(line) || line[i] != '$') {
					segment = segment.WithStop(segment.Start + i - closure)
					if !segment.IsEmpty() {
						node.AppendChild(node, ast.NewRawTextSegment(segment))
					}
					block.Advance(i)
					goto end
				}
			}
		}
		node.AppendChild(node, ast.NewRawTextSegment(segment))
		block.AdvanceLine()
	}
end:
	if !node.IsBlank(block.Source()) {
		// trim first halfspace and last halfspace
		segment := node.FirstChild().(*ast.Text).Segment
		shouldTrimmed := true
		if !(!segment.IsEmpty() && isSpaceOrNewline(block.Source()[segment.Start])) {
			shouldTrimmed = false
		}
		segment = node.LastChild().(*ast.Text).Segment
		if !(!segment.IsEmpty() && isSpaceOrNewline(block.Source()[segment.Stop-1])) {
			shouldTrimmed = false
		}
		if shouldTrimmed {
			t := node.FirstChild().(*ast.Text)
			segment := t.Segment
			t.Segment = segment.WithStart(segment.Start + 1)
			t = node.LastChild().(*ast.Text)
			segment = node.LastChild().(*ast.Text).Segment
			t.Segment = segment.WithStop(segment.Stop - 1)
		}

	}
	return node
}

func isSpaceOrNewline(c byte) bool {
	return c == ' ' || c == '\n'
}

type mathjaxBlockParser struct{}

func NewMathjaxBlockParser() parser.BlockParser {
	return &mathjaxBlockParser{}
}

func (m *mathjaxBlockParser) Trigger() []byte {
	return []byte{'$'}
}

type mathBlockData struct {
	char   byte
	indent int
	length int
	node   ast.Node
}

var mathBlockInfoKey = parser.NewContextKey()

func (m *mathjaxBlockParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	line, segment := reader.PeekLine()
	pos := pc.BlockOffset()
	if pos < 0 || (line[pos] != '$') {
		return nil, parser.NoChildren
	}
	findent := pos
	fenceChar := line[pos]
	i := pos
	for ; i < len(line) && line[i] == fenceChar; i++ {
	}
	oFenceLength := i - pos
	if oFenceLength < 2 {
		return nil, parser.NoChildren
	}
	var info *ast.Text
	if i < len(line)-1 {
		rest := line[i:]
		left := util.TrimLeftSpaceLength(rest)
		right := util.TrimRightSpaceLength(rest)
		if left < len(rest)-right {
			infoStart, infoStop := segment.Start-segment.Padding+i+left, segment.Stop-right
			value := rest[left : len(rest)-right]
			if fenceChar == '$' && bytes.IndexByte(value, '$') > -1 {
				return nil, parser.NoChildren
			} else if infoStart != infoStop {
				info = ast.NewTextSegment(text.NewSegment(infoStart, infoStop))
			}
		}
	}
	node := NewMathBlock(info)
	pc.Set(mathBlockInfoKey, &mathBlockData{fenceChar, findent, oFenceLength, node})
	return node, parser.NoChildren
}

func (m *mathjaxBlockParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	line, segment := reader.PeekLine()
	fdata := pc.Get(mathBlockInfoKey).(*mathBlockData)

	w, pos := util.IndentWidth(line, reader.LineOffset())
	if w < 4 {
		i := pos
		for ; i < len(line) && line[i] == fdata.char; i++ {
		}
		length := i - pos
		if length >= fdata.length && util.IsBlank(line[i:]) {
			newline := 1
			if line[len(line)-1] != '\n' {
				newline = 0
			}
			reader.Advance(segment.Stop - segment.Start - newline + segment.Padding)
			return parser.Close
		}
	}
	pos, padding := util.IndentPositionPadding(line, reader.LineOffset(), segment.Padding, fdata.indent)
	if pos < 0 {
		pos = util.FirstNonSpacePosition(line)
		if pos < 0 {
			pos = 0
		}
		padding = 0
	}
	seg := text.NewSegmentPadding(segment.Start+pos, segment.Stop, padding)
	// if code block line starts with a tab, keep a tab as it is.
	if padding != 0 {
		preserveLeadingTabInCodeBlock(&seg, reader, fdata.indent)
	}
	seg.ForceNewline = true // EOF as newline
	node.Lines().Append(seg)
	reader.AdvanceAndSetPadding(segment.Stop-segment.Start-pos-1, padding)
	return parser.Continue | parser.NoChildren
}

func (m *mathjaxBlockParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {
	fdata := pc.Get(mathBlockInfoKey).(*mathBlockData)
	if fdata.node == node {
		pc.Set(mathBlockInfoKey, nil)
	}
}

func (m *mathjaxBlockParser) CanInterruptParagraph() bool {
	return true
}

func (m *mathjaxBlockParser) CanAcceptIndentedLine() bool {
	return false
}

func preserveLeadingTabInCodeBlock(segment *text.Segment, reader text.Reader, indent int) {
	offsetWithPadding := reader.LineOffset() + indent
	sl, ss := reader.Position()
	reader.SetPosition(sl, text.NewSegment(ss.Start-1, ss.Stop))
	if offsetWithPadding == reader.LineOffset() {
		segment.Padding = 0
		segment.Start--
	}
	reader.SetPosition(sl, ss)
}
