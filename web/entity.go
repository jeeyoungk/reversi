// entities between client & server.
package web

type GameEntity struct {
	ID    int    `json: id`
	Board string `json: string`
}
