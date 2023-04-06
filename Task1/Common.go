package Task1

// Parameters is a struct with parameters for use in a proof of space between a verifier and a prover.
// Id is a random string used to make the proof unique.
// StorageBound is an int which sets a bound on the memory to be used by the prover.
// The storage bound is equal to 2N where N is the amount of nodes in the stored graph
// (that stores bits equal to the hashing size), N should be a power of 2.
// GraphDescription is the graph that should be proven stored, this is specifically a description of the edges.
type Parameters struct {
	Id           string
	StorageBound int
	//TODO Add statistical security parameter (which i used when determining amount of nodes to sample)
	//TODO And the sampling distribution
	GraphDescription Graph
	//TODO Maybe have a different struct for a pebbled graph (that indcludes values), and one that is only for edges and size
}
