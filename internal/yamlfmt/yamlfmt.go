package yamlfmt

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/yaml.v3"
	"io"
	"strings"
)

var _ = spew.Dump // prevent spew from being removed if unused

// Format takes a yaml bytearray and applies the relevant changes to format it
func Format(in []byte) ([]byte, error) {

	yamlContent := in

	var err error
	yamlContent, err = addLineBreakFlag(yamlContent)
	if err != nil {
		return nil, err
	}

	// Format the yaml content
	reader := bytes.NewReader(yamlContent)
	decoder := yaml.NewDecoder(reader)
	documents := []yaml.Node{}

	for {
		var docNode yaml.Node
		err := decoder.Decode(&docNode)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("error decoding yaml document: %v", err)
		}
		walkNode(&docNode)
		documents = append(documents, docNode)
	}

	var b bytes.Buffer
	e := yaml.NewEncoder(&b)
	e.SetIndent(2)
	for _, doc := range documents {
		err = e.Encode(&doc)
		if err != nil {
			return nil, err
		}
	}

	result := b.Bytes()
	result, err = removeLineBreakFlag(result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

const magicLine = "#==================xieK7ej3eiwoh9mie4che=========================="

// addLineBreakFlag reads all the lines in the file and replaces any line breaks with a single flag to be removed later
func addLineBreakFlag(in []byte) ([]byte, error) {

	reader := bytes.NewReader(in)
	scanner := bufio.NewScanner(reader)

	out := strings.Builder{}
	lineCount := 0
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSuffix(line, " ")

		if strings.TrimSpace(line) == "" {
			if lineCount == 0 {
				lineCount = 1
				out.WriteString(magicLine + "\n")
			}
		} else {
			lineCount = 0
			out.WriteString(line + "\n")
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error when readeing a line with scanner: %v", err)
	}

	return []byte(out.String()), nil
}

// removeLineBreakFlag reads all the lines in the file and replaces linebreak flags with real line breaks
// sometimes line breaks are still present in the output, this function will remove them as well and only
// introduce line breaks where it finds the magic line
func removeLineBreakFlag(in []byte) ([]byte, error) {
	reader := bytes.NewReader(in)
	scanner := bufio.NewScanner(reader)

	out := strings.Builder{}

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSuffix(line, " ")
		if strings.TrimSpace(line) == "" {
			continue
		}
		if strings.TrimSpace(line) == magicLine {
			out.WriteString("\n")
		} else {
			out.WriteString(line + "\n")
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error when readeing a line with scanner: %v", err)
	}
	return []byte(out.String()), nil
}

// walkNode will recursively walk the yaml structure and apply changes:
// - remove quote style from keys and values
func walkNode(node *yaml.Node) {

	if len(node.Content) > 0 {
		for _, n := range node.Content {
			if n.Style == yaml.DoubleQuotedStyle {
				n.Style = 0
			}
			walkNode(n)
		}
	}
}
