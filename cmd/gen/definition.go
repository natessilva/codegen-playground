package main

type Definition struct {
	Services []Service
	Structs  []Struct
}

type Service struct {
	Name    string
	Methods []Method
}

type Method struct {
	Name   string
	Input  string
	Output string
}

type Struct struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name      string
	IsSlice   bool
	IsNumeric bool
	Type      string
	Tag       string
}
