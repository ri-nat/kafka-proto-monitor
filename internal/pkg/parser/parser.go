package parser

import (
	"fmt"
	"path/filepath"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/ri-nat/kafka-proto-monitor/internal/pkg/config"
)

type Parser struct {
	File *desc.FileDescriptor
}

func NewParser(opts *config.Options) (*Parser, error) {
	// Set folder that contains proto file as import path
	//
	// TODO: add ability to configure import paths via cli arguments
	dir := filepath.Dir(opts.ProtoFile)
	importPaths := []string{dir}

	// Resolve and parse proto
	fileNames, err := protoparse.ResolveFilenames(importPaths, opts.ProtoFile)
	if err != nil {
		return nil, err
	}

	p := protoparse.Parser{ImportPaths: importPaths}
	linkedFiles, err := p.ParseFiles(fileNames...)
	if err != nil {
		return nil, err
	}

	parser := Parser{
		File: linkedFiles[0],
	}

	// Look for a specific message in a parsed proto, if it is specified in options
	if len(opts.ProtoMsg) > 0 {
		if _, err := parser.findProtoMessage(opts.ProtoMsg); err != nil {
			return nil, err
		}
		fmt.Println("Proto message found")
	} else {
		fmt.Println("Working with all messages in proto file")
	}

	return &parser, nil
}

func (p *Parser) NewMessageFromBytes(protoName string, buf []byte) (*dynamic.Message, error) {
	dm, err := p.NewMessage(protoName)
	if err != nil {
		return nil, err
	}

	if err := dm.Unmarshal(buf); err != nil {
		return nil, fmt.Errorf("can't unmarshal bytes to message: %w", err)
	}

	return dm, nil
}

func (p *Parser) NewMessage(protoName string) (*dynamic.Message, error) {
	msg, err := p.findProtoMessage(protoName)
	if err != nil {
		return nil, err
	}

	return dynamic.NewMessage(msg), nil
}

func (p *Parser) findProtoMessage(protoName string) (*desc.MessageDescriptor, error) {
	md := p.File.FindMessage(protoName)
	if md == nil {
		return nil, fmt.Errorf("can't find proto message `%s` in proto file", protoName)
	}

	return md, nil
}
