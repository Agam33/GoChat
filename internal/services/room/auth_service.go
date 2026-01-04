package room

type RoomService interface {
	CreateRoom()
	GetRoomById()
}

type roomService struct {
	roomRepo RoomRepository
}

func NewRoomService(roomRepo RoomRepository) RoomService {
	return &roomService{
		roomRepo: roomRepo,
	}
}

func (r *roomService) CreateRoom() {}

func (r *roomService) GetRoomById() {}
