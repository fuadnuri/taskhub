package service 

type UserService struct{
	repo UserRepository
}

func NewUserService(repo UserRepository)*UserService{
	return &UserService{
		repo:repo,
	}
}

func (s *UserService) CreateUser(user *entities.UserEntity) error{
	
	return s.repo.CreateUser(user)
}


func (s *UserService) GetUserById(id uint)(*entities.UserEntity,error){
	//something goes here later for now 
	return s.repo.GetUserById(id)
}

func (s *UserService)GetAllUsers()([]entities.UserEntity,error){
	return s.repo.GetAllUsers()
}