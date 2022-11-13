package campaign

type Service interface {
	GetCampaigns(userID int) ([]Campaign, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetCampaigns(userID int) ([]Campaign, error) {

	//kalau ada ID usernya, maka cari campaign yang ada user ID bersangkutan
	if userID != 0 {
		campaigns, err := s.repository.FindByUserID(userID)

		if err != nil {
			return campaigns, err
		}

		return campaigns, nil
	}

	//kalau tidak ada, tampilkan semua
	campaigns, err := s.repository.FindAll()

	if err != nil {
		return campaigns, err
	}

	return campaigns, nil
}
