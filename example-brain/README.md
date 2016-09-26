# Example Brain

This small Go program gives an example of how one could implement a brain (a server that provides a scheduling algorithm to Layer-X)

To run, just:

```bash
go build
./example-brain -core <ip:port>
``` 

and watch your tasks automatically get assigned in round-robin fashion