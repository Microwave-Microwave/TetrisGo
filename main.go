package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// -----------------------------------------------------------------------//
// --------------------------   Structs   --------------------------------//
// -----------------------------------------------------------------------//
type Tuple struct {
	x, y int32
}
type Coordinate = Tuple
type Vector = Tuple

type Parameters struct {
	grid_size, target_fps int32
	play_area             Tuple
	screen_size           Coordinate
	window_name           string
}

type Cell struct {
	entity_type string
	color       rl.Color
	style       string
}

type Board struct {
	grid [][]Cell
}

type Game struct {
	pars   Parameters
	board  Board
	pieces []Piece
}

// -----------------------------------------------------------------------//
// ----------------------   Pieces and Player   --------------------------//
// -----------------------------------------------------------------------//
type Piece struct {
	name   string
	center Coordinate
	count  int
	pieces [4]Coordinate
}

func create_piece(name string, position Coordinate) (p Piece) {
	p.name = name
	p.center = position
	p.count = 4
	x, y := position.x, position.y
	switch name {
	case "o":
		p.pieces[0] = Coordinate{x, y}
		p.pieces[1] = Coordinate{x + 1, y}
		p.pieces[2] = Coordinate{x, y + 1}
		p.pieces[3] = Coordinate{x + 1, y + 1}
	case "i":
		p.pieces[0] = Coordinate{x, y}
		p.pieces[1] = Coordinate{x + 1, y}
		p.pieces[2] = Coordinate{x + 2, y}
		p.pieces[3] = Coordinate{x + 3, y}
	case "s":
		p.pieces[0] = Coordinate{x, y}
		p.pieces[1] = Coordinate{x + 1, y}
		p.pieces[2] = Coordinate{x + 1, y + 1}
		p.pieces[3] = Coordinate{x + 2, y + 1}
	case "z":
		p.pieces[0] = Coordinate{x, y + 1}
		p.pieces[1] = Coordinate{x + 1, y + 1}
		p.pieces[2] = Coordinate{x + 1, y}
		p.pieces[3] = Coordinate{x + 2, y}
	case "l":
		p.pieces[0] = Coordinate{x, y}
		p.pieces[1] = Coordinate{x, y}
		p.pieces[2] = Coordinate{x, y}
		p.pieces[3] = Coordinate{x, y}
	case "j":
		p.pieces[0] = Coordinate{x, y}
		p.pieces[1] = Coordinate{x, y + 1}
		p.pieces[2] = Coordinate{x, y + 2}
		p.pieces[3] = Coordinate{x + 1, y}
	case "t":
		p.pieces[0] = Coordinate{x, y + 1}
		p.pieces[1] = Coordinate{x + 1, y + 1}
		p.pieces[2] = Coordinate{x + 2, y + 1}
		p.pieces[3] = Coordinate{x + 1, y}
	default:
		fmt.Println("ERROR")
	}
	return p
}

type Matrix[T any] struct {
	w, h int
	data []T
}

func rotate_matrix(matrix [6][6]int) {
	var rotated_matrix [6][6]int

	for y := 0; y < 6; y++ {
		for x := 0; x < 6; x++ {
			rotated_matrix[x][y] = matrix[x][y]
		}
	}

	for x in range(0, 6):
	for y in range(0, 6):
	    #2x right
	    rotated_matrix[x][y] = matrix[5-x][5-y]
	    #right
	    rotated_matrix[x][y] = matrix[5-y][x]
	    #left
	    rotated_matrix[5-y][x] = matrix[x][y]

}

func rotate_piece(piece *Piece) {
	//6x6 empty matrix
	var matrix [6][6]int
	for x := 0; x < 6; x++ {
		for y := 0; y < 6; y++ {
			matrix[x][y] = 0
		}
	}

	origin := Coordinate{piece.pieces[0].x, piece.pieces[0].y}
	new_center := origin
	new_center.x += 1
	new_center.y += 1
	matrix[new_center.x][new_center.y] = 1
	for i := 1; i < 4; i++ {
		cell := Coordinate{piece.pieces[i].x, piece.pieces[i].y}
		matrix[cell.x-origin.x][cell.y-origin.y] = 1
	}

}

func fix_rotation(game *Game, piece *Piece) {
	//TODO
}

func spawn_piece(game *Game, name string, position Coordinate, rotation int) {
	var piece = create_piece(name, position)
	rotate_piece(&piece)
	fix_rotation(game, &piece)
}

type Player struct {
	piece Piece
}

// -----------------------------------------------------------------------//
func new_cell(entity_type string, color rl.Color) Cell {
	return Cell{
		entity_type: entity_type,
		color:       color,
		style:       "borderless",
	}
}

func new_cell_with_style(entity_type string, color rl.Color, style string) Cell {
	return Cell{
		entity_type: entity_type,
		color:       color,
		style:       style,
	}
}

// -----------------------------------------------------------------------//
// --------------------------   Functions   ------------------------------//
// -----------------------------------------------------------------------//

func (c Cell) print() {
	var cell_text string

	switch c.entity_type {
	case "background":
		cell_text = "0"
	case "border":
		cell_text = "1"
	case "player":
		cell_text = "2"
	case "junk":
		cell_text = "E"
	}

	fmt.Print(cell_text, " ")

}

func (b Board) print() {
	fmt.Println("Matrix")

	for y := range b.grid[0] {
		for x := range b.grid {
			b.grid[x][y].print()
		}
		fmt.Println()
	}
	fmt.Println()

}

func create_game() Game {
	var parameters Parameters = parameter_setup()
	var board Board = board_setup(parameters)
	var pieces []Piece = []Piece{} //TODO
	return Game{parameters, board, pieces}
}

func parameter_setup() (par Parameters) {
	par.grid_size = 32
	par.play_area = Tuple{12, 24}
	par.screen_size = Coordinate{par.grid_size * par.play_area.x, par.grid_size * par.play_area.y} // 12 * grid | 24 * grid
	par.window_name = "Tetris"
	par.target_fps = 60
	return
}

func traverse(start, end, steps int, fn func(int)) {
	for i := start; i < end; i += steps {
		fn(i)
	}
}

func set_board_style(board Board, style string) Board {
	var x_count = len(board.grid)
	var y_count = len(board.grid[0])
	traverse(0, x_count, 1, func(x int) {
		traverse(0, y_count, 1, func(y int) {
			board.grid[x][y].style = style
		})
	})
	return board
}

func board_setup(par Parameters) (board Board) {
	var x_count = int(par.play_area.x)
	var y_count = int(par.play_area.y)
	grid := make([][]Cell, x_count)

	white_cell := new_cell("background", rl.White)
	black_cell := new_cell("background", rl.Black)

	// create matrix and paint background white
	traverse(0, x_count, 1, func(x int) {
		grid[x] = make([]Cell, y_count)
		traverse(0, y_count, 1, func(y int) {
			grid[x][y] = white_cell
		})
	})

	// paint the two sides black
	traverse(0, y_count, 1, func(n int) {
		grid[0][n] = black_cell
		grid[x_count-1][n] = black_cell
	})

	// paint the bottom black
	traverse(0, x_count, 1, func(n int) {
		grid[n][y_count-1] = black_cell
	})
	board.grid = grid
	board = set_board_style(board, "bordered")

	return
}

func window_setup(game *Game) {
	par := game.pars
	var x, y int32 = par.screen_size.x, par.screen_size.y
	rl.InitWindow(x, y, par.window_name)
	rl.SetTargetFPS(int32(par.target_fps))
}

func draw_square(start Coordinate, size int32, color rl.Color) {
	rl.DrawRectangle(start.x, start.y, size, size, color)
}

func draw_cell_with_border(start Coordinate, grid_size int32, inner_color rl.Color, border_size int32, border_color rl.Color) {
	border_cell_start := Coordinate{
		start.x * grid_size,
		start.y * grid_size,
	}
	draw_square(border_cell_start, grid_size, border_color)

	inner_cell_start := Coordinate{
		start.x*grid_size + border_size,
		start.y*grid_size + border_size,
	}
	innerl_cell_size := grid_size - border_size*2
	draw_square(inner_cell_start, innerl_cell_size, inner_color)
}

func draw_cell_without_border(start Coordinate, grid_size int32, color rl.Color) {
	border_cell_start := Coordinate{
		start.x * grid_size,
		start.y * grid_size,
	}
	draw_square(border_cell_start, grid_size, color)
}

func draw_cell(game *Game, x int, y int, style string, layer int) {
	start := Coordinate{int32(x), int32(y)}
	cell := game.board.grid[x][y]
	switch style {
	//TODO MAKE LAYERS DIFFERENT MAYBE
	case "borderless":
		if layer == 0 {
			draw_cell_without_border(start, game.pars.grid_size, cell.color)
		} else if layer == 1 {
			draw_cell_without_border(start, game.pars.grid_size, cell.color)
		} else if layer == 2 {
			draw_cell_without_border(start, game.pars.grid_size, cell.color)
		}
	case "bordered":
		draw_cell_with_border(start, game.pars.grid_size, game.board.grid[x][y].color, 1, rl.Gray)
	}

}

func draw_board(game *Game, layer int) {
	rl.BeginDrawing()
	rl.ClearBackground(rl.RayWhite)

	var x_count = int(game.pars.play_area.x)
	var y_count = int(game.pars.play_area.y)
	traverse(0, x_count, 1, func(x int) {
		traverse(0, y_count, 1, func(y int) {
			var style = game.board.grid[x][y].style
			draw_cell(game, x, y, style, layer)
		})
	})

	rl.EndDrawing()
}

func draw_arena(game *Game) {
	draw_board(game, 0)
}

func draw_entities(game *Game) {
	draw_board(game, 1)
}

func draw_player(game *Game) {
	draw_board(game, 2)
}

func game_loop(game *Game) {

	for !rl.WindowShouldClose() {
		draw_arena(game)
		draw_entities(game)
		draw_player(game)
	}
}

// -----------------------------------------------------------------------//
// ---------------------------   Main   ----------------------------------//
// -----------------------------------------------------------------------//
func main() {
	var game = create_game()
	game.board.print()

	window_setup(&game)
	game_loop(&game)
	defer rl.CloseWindow()
}
