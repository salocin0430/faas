package entity

type Container struct {
	ID        string
	ImageURL  string
	Status    string
	Resources ResourceLimits
}

type ResourceLimits struct {
	Memory    int64
	CPUShares int64
}
