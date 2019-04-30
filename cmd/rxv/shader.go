package main

// Vertex Shader
// This is pass-through vertex shader which
// sends its input directly to the geometry shader
// without any processing.
//
const sourceGSDemoVertex = `
#include <attributes>

// Outputs for geometry shader
out vec3 vnormal;

void main() {

	gl_Position = vec4(VertexPosition, 1.0);
  	vnormal = VertexNormal;
}

`

//
// Geometry Shader
// This geometry shader receives triangles vertices
// from the vertex shader and generates lines for
// wireframe and/or vertex normals and/or face normals.
//
const sourceGSDemoGeometry = `
layout (triangles) in;
layout (line_strip, max_vertices = 12) out;

// Model uniforms
uniform mat4 MVP;

// Inputs from Vertex Shader
in vec3 vnormal[];

// Inputs uniforms
int ShowWireframe=1;
int ShowVnormal=1;
int ShowFnormal=0;

// Colors
const vec4 colorWire    = vec4(1, 1, 0, 1);
const vec4 colorVnormal = vec4(1, 0, 0, 1);
const vec4 colorFnormal = vec4(0, 0, 1, 1);

// Output color to fragment shader
out vec4 vertex_color;

void main() {

	// Emits triangle's vertices as lines to show wireframe
	if (ShowWireframe != 0) {
		for (int n = 0; n < gl_in.length(); n++) {
			// Vertex position
			gl_Position = MVP * gl_in[n].gl_Position;
			vertex_color = colorWire;
			EmitVertex();
		}
		// Emit first triangle vertex to close the last line strip.
		gl_Position = MVP * gl_in[0].gl_Position;
		vertex_color = colorWire;
		EmitVertex();
		EndPrimitive();
	}

	// Emits lines representing the vertices normals
	if (ShowVnormal != 0) {
		for (int i = 0; i < gl_in.length(); i++) {

			vec3 position = gl_in[i].gl_Position.xyz;
			vec3 normal = vnormal[i];

			gl_Position = MVP * vec4(position, 1.0);
			vertex_color = colorVnormal;
			EmitVertex();

			gl_Position = MVP * vec4(position + normal * 0.5, 1.0);
			vertex_color = colorVnormal;
			EmitVertex();

			EndPrimitive();
		}
	}

	// Emits one line representing the face normal
	if (ShowFnormal != 0) {
		vec3 p0 = gl_in[0].gl_Position.xyz;
		vec3 p1 = gl_in[1].gl_Position.xyz;
		vec3 p2 = gl_in[2].gl_Position.xyz;

		vec3 v0 = p0 - p1;
		vec3 v1 = p2 - p1;
		vec3 faceN = normalize(cross(v1, v0));

		// Center of the triangle
		vec3 center = (p0 + p1 + p2) / 3.0;

		gl_Position = MVP * vec4(center, 1.0);
		vertex_color = colorFnormal;
		EmitVertex();

		gl_Position = MVP * vec4(center + faceN * 0.5, 1.0);
		vertex_color = colorFnormal;
		EmitVertex();
		EndPrimitive();
	}
}

`

//
// Fragment Shader template
//
const sourceGSDemoFrag = `
in vec4 vertex_color;
out vec4 Out_Color;

void main() {
	Out_Color = vertex_color;
}

`
