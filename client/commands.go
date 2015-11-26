package client

import (
	"github.com/nkatsaros/pipboygo/protocol"
)

type Response struct {
	Allowed bool
	Success bool
	Message string
}

func (c *Client) UseItem(handleid, version int) error {
	_, err := c.enc.Encode(protocol.Command{0, []interface{}{handleid, 0, version}})
	return err
}

func (c *Client) DropItem(handleid, count, version int, stackid []int) error {
	_, err := c.enc.Encode(protocol.Command{1, []interface{}{handleid, count, version, stackid}})
	return err
}

func (c *Client) FavoriteItem(handleid, position, version int, stackid []int) error {
	_, err := c.enc.Encode(protocol.Command{2, []interface{}{handleid, 0, version}})
	return err
}

func (c *Client) TagForSearch(componentformid, version int) error {
	_, err := c.enc.Encode(protocol.Command{3, []interface{}{componentformid, version}})
	return err
}

func (c *Client) CycleSearch(page int) error {
	_, err := c.enc.Encode(protocol.Command{4, []interface{}{page}})
	return err
}

func (c *Client) ToggleQuestMarker(questid int) error {
	_, err := c.enc.Encode(protocol.Command{5, []interface{}{questid, 0, 0}})
	return err
}

func (c *Client) PlaceCustomMarker(x, y float32) error {
	_, err := c.enc.Encode(protocol.Command{6, []interface{}{x, y, false}})
	return err
}

func (c *Client) RemoveCustomMarker() error {
	_, err := c.enc.Encode(protocol.Command{7, []interface{}{}})
	return err
}

func (c *Client) FastTravel(id int) (Response, error) {
	_, err := c.enc.Encode(protocol.Command{9, []interface{}{id}})
	return Response{}, err
}

func (c *Client) ToggleRadio(id int) error {
	_, err := c.enc.Encode(protocol.Command{12, []interface{}{id}})
	return err
}

func (c *Client) RequestLocalMapUpdate() error {
	_, err := c.enc.Encode(protocol.Command{13, []interface{}{}})
	return err
}

func (c *Client) Refresh() error {
	_, err := c.enc.Encode(protocol.Command{14, []interface{}{}})
	return err
}
