package rclient

import "fmt"

func (rc *RClient) ListProjects() error {

	if rc.SOAP == nil {
		fmt.Println("Client API is not initializated")
		return nil
	}

	p, e := rc.API.ListProjects()
	fmt.Println(e)
	fmt.Println(p)
	for _, v := range p {
		fmt.Println(v.Name)
	}
	return nil
}
