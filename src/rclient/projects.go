package rclient

import "fmt"

func (rc *RClient) ListProjects() error {
	p, e := rc.API.ListProjects()
	fmt.Println(e)
	fmt.Println(p)
	for _, v := range p {
		fmt.Println(v.Name)
	}
	return nil
}
