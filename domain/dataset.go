package domain

type Dataset struct {
	Id string

	Owner    Account
	Protocol ProtocolName

	DatasetModifiableProperty

	RepoId string

	CreatedAt int64
	UpdatedAt int64

	Version int

	// following fileds is not under the controlling of version
	LikeCount     int
	DownloadCount int

	RelatedModels   RelatedResources
	RelatedProjects RelatedResources
}

func (d *Dataset) IsPrivate() bool {
	return d.RepoType.RepoType() == RepoTypePrivate
}

type DatasetModifiableProperty struct {
	Name     DatasetName
	Desc     ResourceDesc
	RepoType RepoType
	Tags     []string
}

type DatasetSummary struct {
	Id            string
	Owner         Account
	Name          DatasetName
	Desc          ResourceDesc
	Tags          []string
	UpdatedAt     int64
	LikeCount     int
	DownloadCount int
}
