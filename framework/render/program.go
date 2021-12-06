package render

import (
	"fmt"
	gl "github.com/go-gl/gl/v3.1/gles2"
	"strings"
)

type Attribute struct {
	glType   uint32
	size     int32
	rows     uint32
	buffer   uint32
	location uint32
}

type Uniform struct {
	glType   uint32
	location int32
}

type Program struct {
	program       uint32
	attributes    map[string]Attribute
	uniforms      map[string]Uniform
	texture_units [8]uint32
	draw_size     int32
}

func (p *Program) loadShaders(vert_file string, frag_file string) {
	//load & compile the shaders
	vert_shader := gl.CreateShader(gl.VERTEX_SHADER)
	vert_source, err := Resources.ReadFile("resources/shaders/" + vert_file)
	if err != nil {
		panic(err)
	}
	vert_str, free := gl.Strs(string(vert_source) + "\x00")
	gl.ShaderSource(vert_shader, 1, vert_str, nil)
	free()
	gl.CompileShader(vert_shader)
	vert_compile_status := int32(0)
	gl.GetShaderiv(vert_shader, gl.COMPILE_STATUS, &vert_compile_status)
	if vert_compile_status == 0 {
		fmt.Println("VERTEX SHADER ERROR:")
		var length int32
		gl.GetShaderiv(vert_shader, gl.INFO_LOG_LENGTH, &length)
		buf := strings.Repeat("\x00", int(length+1))
		gl.GetShaderInfoLog(vert_shader, length, nil, gl.Str(buf))
		fmt.Printf("%s\n", buf)
	}
	frag_shader := gl.CreateShader(gl.FRAGMENT_SHADER)
	frag_source, err := Resources.ReadFile("resources/shaders/" + frag_file)
	if err != nil {
		panic(err)
	}
	frag_str, free := gl.Strs(string(frag_source) + "\x00")
	gl.ShaderSource(frag_shader, 1, frag_str, nil)
	free()
	gl.CompileShader(frag_shader)
	frag_compile_status := int32(0)
	gl.GetShaderiv(frag_shader, gl.COMPILE_STATUS, &frag_compile_status)
	if frag_compile_status == 0 {
		fmt.Println("FRAGMENT SHADER ERROR:")
		var length int32
		gl.GetShaderiv(frag_shader, gl.INFO_LOG_LENGTH, &length)
		buf := strings.Repeat("\x00", int(length+1))
		gl.GetShaderInfoLog(frag_shader, length, nil, gl.Str(buf))
		fmt.Printf("%s\n", buf)
	}
	program := gl.CreateProgram()
	gl.AttachShader(program, vert_shader)
	gl.AttachShader(program, frag_shader)
	gl.LinkProgram(program)
	link_status := int32(0)
	gl.GetProgramiv(program, gl.LINK_STATUS, &link_status)
	if link_status == 0 {
		fmt.Println("SHADER PROGRAM ERROR:")
		var length int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &length)
		buf := strings.Repeat("\x00", int(length+1))
		gl.GetProgramInfoLog(program, length, nil, gl.Str(buf))
		fmt.Printf("%s\n", buf)
	}
	//fmt.Printf("program id: %d\n", program)

	p.program = program
}

func (p *Program) scanAttributes() {
	//scan through active attributes & create buffers & layouts for them
	var numAttributes int32
	gl.GetProgramiv(p.program, gl.ACTIVE_ATTRIBUTES, &numAttributes)
	fmt.Printf("numAttributes: %d\n", numAttributes)
	for i := 0; i < int(numAttributes); i++ {
		var buf [64]uint8
		var glType uint32
		var arraySize int32
		gl.GetActiveAttrib(p.program, uint32(i), int32(len(buf)), nil, &arraySize, &glType, &buf[0])
		name := gl.GoStr(&buf[0])
		location := uint32(gl.GetAttribLocation(p.program, &buf[0]))
		size := int32(0)
		rows := uint32(0)
		buffer := uint32(0)
		gl.GenBuffers(1, &buffer)
		switch glType {
		case gl.FLOAT:
			size = 1
			rows = 1
		case gl.FLOAT_VEC2:
			size = 2
			rows = 1
		case gl.FLOAT_VEC3:
			size = 3
			rows = 1
		case gl.FLOAT_VEC4:
			size = 4
			rows = 1
		case gl.FLOAT_MAT2:
			size = 2
			rows = 2
		case gl.FLOAT_MAT3:
			size = 3
			rows = 3
		case gl.FLOAT_MAT4:
			size = 4
			rows = 4
		}
		p.attributes[name] = Attribute{glType, size, rows, buffer, location}
		fmt.Printf("attribute %d (%dx%d): %s\n", location, rows, size, name)
	}
}

func (p *Program) scanUniforms() {
	var numUniforms int32
	gl.GetProgramiv(p.program, gl.ACTIVE_UNIFORMS, &numUniforms)
	fmt.Printf("numUniforms: %d\n", numUniforms)
	for i := 0; i < int(numUniforms); i++ {
		var buf [64]uint8
		var arraySize int32
		var glType uint32
		gl.GetActiveUniform(p.program, uint32(i), int32(len(buf)), nil, &arraySize, &glType, &buf[0])
		name := gl.GoStr(&buf[0])
		location := gl.GetUniformLocation(p.program, &buf[0])
		p.uniforms[name] = Uniform{glType, location}
		fmt.Printf("uniform %d: %s\n", location, name)
	}
}

func CreateProgram(vert_source string, frag_source string) *Program {
	p := new(Program)
	p.loadShaders(vert_source, frag_source)
	p.attributes = make(map[string]Attribute)
	p.scanAttributes()
	p.uniforms = make(map[string]Uniform)
	p.scanUniforms()
	return p
}

func (p *Program) BufferData(name string, data []float32) {
	if attr, ok := p.attributes[name]; ok {
		gl.BindBuffer(gl.ARRAY_BUFFER, attr.buffer)
		gl.BufferData(gl.ARRAY_BUFFER, len(data)*4, gl.Ptr(data), gl.STATIC_DRAW)
		p.draw_size = int32(len(data)) / (attr.size * int32(attr.rows))
	}
}

//for long-lived quick-swap storage i.e. for models & such
type Buffer struct {
	Data   []float32
	buffer uint32
}

func (p *Program) BindBuffer(name string, buffer *Buffer) {
	if buffer.buffer == 0 {
		//send buffer to GPU if not already
		gl.GenBuffers(1, &buffer.buffer)
		gl.BindBuffer(gl.ARRAY_BUFFER, buffer.buffer)
		gl.BufferData(gl.ARRAY_BUFFER, len(buffer.Data)*4, gl.Ptr(buffer.Data), gl.STATIC_DRAW)
	}
	if attr, ok := p.attributes[name]; ok {
		//replace the default attribute buffer
		attr.buffer = buffer.buffer
		p.attributes[name] = attr
		p.draw_size = int32(len(buffer.Data)) / (attr.size * int32(attr.rows))
	}
}

func (p *Program) Uniform(name string, value interface{}) {
	gl.UseProgram(p.program)
	if unif, ok := p.uniforms[name]; ok {
		switch unif.glType {
		case gl.FLOAT:
			gl.Uniform1f(unif.location, value.(float32))
		case gl.FLOAT_VEC2:
			v := value.([2]float32)
			gl.Uniform2fv(unif.location, 1, &v[0])
		case gl.FLOAT_VEC3:
			v := value.([3]float32)
			gl.Uniform3fv(unif.location, 1, &v[0])
		case gl.FLOAT_VEC4:
			v := value.([4]float32)
			gl.Uniform4fv(unif.location, 1, &v[0])
		case gl.FLOAT_MAT2:
			v := value.([2 * 2]float32)
			gl.UniformMatrix2fv(unif.location, 1, false, &v[0])
		case gl.FLOAT_MAT3:
			v := value.([3 * 3]float32)
			gl.UniformMatrix3fv(unif.location, 1, false, &v[0])
		case gl.FLOAT_MAT4:
			v := value.([4 * 4]float32)
			gl.UniformMatrix4fv(unif.location, 1, false, &v[0])
		case gl.INT:
			gl.Uniform1i(unif.location, value.(int32))
		case gl.INT_VEC2:
			v := value.([2]int32)
			gl.Uniform2iv(unif.location, 1, &v[0])
		case gl.INT_VEC3:
			v := value.([3]int32)
			gl.Uniform3iv(unif.location, 1, &v[0])
		case gl.INT_VEC4:
			v := value.([4]int32)
			gl.Uniform4iv(unif.location, 1, &v[0])
		case gl.BOOL:
			gl.Uniform1i(unif.location, value.(int32))
		case gl.BOOL_VEC2:
			v := value.([2]int32)
			gl.Uniform2iv(unif.location, 1, &v[0])
		case gl.BOOL_VEC3:
			v := value.([3]int32)
			gl.Uniform3iv(unif.location, 1, &v[0])
		case gl.BOOL_VEC4:
			v := value.([4]int32)
			gl.Uniform4iv(unif.location, 1, &v[0])
		case gl.SAMPLER_2D:
			for free_unit := 0; free_unit < 8; free_unit++ {
				if p.texture_units[free_unit] == value {
					break
				}
				if p.texture_units[free_unit] == 0 {
					p.texture_units[free_unit] = value.(uint32)
					gl.ActiveTexture(uint32(gl.TEXTURE0 + free_unit))
					gl.BindTexture(gl.TEXTURE_2D, value.(uint32))
					gl.Uniform1i(unif.location, int32(free_unit))
				}
			}
		case gl.SAMPLER_CUBE:
			for free_unit := 0; free_unit < 8; free_unit++ {
				if p.texture_units[free_unit] == 0 {
					break
				}
				if p.texture_units[free_unit] == 0 {
					p.texture_units[free_unit] = value.(uint32)
					gl.ActiveTexture(uint32(gl.TEXTURE0 + free_unit))
					gl.BindTexture(gl.TEXTURE_CUBE_MAP, value.(uint32))
					gl.Uniform1i(unif.location, int32(free_unit))
				}
			}
		}
	}
}

func (p *Program) LoadAttributes() {
	gl.UseProgram(p.program)
	for _, attr := range p.attributes {
		gl.BindBuffer(gl.ARRAY_BUFFER, attr.buffer)
		for row := uint32(0); row < attr.rows; row++ {
			gl.EnableVertexAttribArray(attr.location + row)
			row_size := 4 * attr.size //size of the individual row (i.e. the vecn)
			gl.VertexAttribPointerWithOffset(attr.location+row, attr.size, gl.FLOAT, false, row_size*int32(attr.rows), uintptr(row_size*int32(row)))
		}
	}
}
func (p *Program) DrawArrays() {
	//basically just the bare draw call without the setup/teardown, so it can be repeated easily.
	gl.DrawArrays(gl.TRIANGLES, 0, p.draw_size)
}
func (p *Program) ClearTextureUnits() {
	for i := 0; i < 8; i++ {
		p.texture_units[i] = 0
	}
}

func (p *Program) Draw() {
	p.LoadAttributes()
	p.DrawArrays()
	p.ClearTextureUnits()
}
