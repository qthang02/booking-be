package database

import (
	"github.com/qthang02/booking/enities"
	"github.com/qthang02/booking/types"
	"github.com/qthang02/booking/util"
	"math/rand"
	"time"
)

func InitUsersDataDefault() []*enities.User {
	data := []*enities.User{
		{Name: "Nguyen Quoc Thang", Username: "thangnq", Email: "thangnq1@vnpay.vn", Phone: "0947013683", Birthday: randomBirthday(), Gender: true, Address: "Quang Trung, Go Vap", Password: "admin@123456", Role: util.Admin},
		{Name: "Le Trung Tin", Username: "tinle", Email: "lenguyen170804@gmail.com", Phone: "098113579", Birthday: randomBirthday(), Gender: true, Address: "Tan Binh", Password: "admin@123456", Role: util.Admin},
		{Name: "Alice Johnson", Username: "alicej", Email: "alice@example.com", Phone: "123-456-7890", Birthday: randomBirthday(), Gender: true, Address: "123 Main St", Password: "123456", Role: util.Customer},
		{Name: "Bob Smith", Username: "bobsmith", Email: "bob@example.com", Phone: "234-567-8901", Birthday: randomBirthday(), Gender: true, Address: "456 Elm St", Password: "123456", Role: util.Customer},
		{Name: "Carol White", Username: "carolw", Email: "carol@example.com", Phone: "345-678-9012", Birthday: randomBirthday(), Gender: false, Address: "789 Oak St", Password: "123456", Role: util.Customer},
		{Name: "David Brown", Username: "davidb", Email: "david@example.com", Phone: "456-789-0123", Birthday: randomBirthday(), Gender: true, Address: "101 Pine St", Password: "123456", Role: util.Customer},
		{Name: "Eve Black", Username: "eveb", Email: "eve@example.com", Phone: "567-890-1234", Birthday: randomBirthday(), Gender: false, Address: "202 Maple St", Password: "123456", Role: util.Customer},
		{Name: "Frank Green", Username: "frankg", Email: "frank@example.com", Phone: "678-901-2345", Birthday: randomBirthday(), Gender: true, Address: "303 Birch St", Password: "123456", Role: util.Customer},
		{Name: "Grace Red", Username: "gracer", Email: "grace@example.com", Phone: "789-012-3456", Birthday: randomBirthday(), Gender: false, Address: "404 Cedar St", Password: "123456", Role: util.Customer},
		{Name: "Henry Blue", Username: "henryb", Email: "henry@example.com", Phone: "890-123-4567", Birthday: randomBirthday(), Gender: true, Address: "505 Walnut St", Password: "123456", Role: util.Customer},
		{Name: "Ivy Purple", Username: "ivyp", Email: "ivy@example.com", Phone: "901-234-5678", Birthday: randomBirthday(), Gender: false, Address: "606 Spruce St", Password: "123456", Role: util.Customer},
		{Name: "Jack Gold", Username: "jackg", Email: "jack@example.com", Phone: "012-345-6789", Birthday: randomBirthday(), Gender: true, Address: "707 Fir St", Password: "123456", Role: util.Customer},
	}

	return data
}

func InitOrdersDataDefault() []*enities.Order {

	data := []*enities.Order{
		{GuestNumber: 2, Price: 100.5, Description: "Order 1", Checkin: time.Now(), Checkout: time.Now().Add(48 * time.Hour), CategoryType: types.VIP1, RoomNumber: 101, UserID: 1},
		{GuestNumber: 4, Price: 200.75, Description: "Order 2", Checkin: time.Now(), Checkout: time.Now().Add(72 * time.Hour), CategoryType: types.DELUXE, RoomNumber: 102, UserID: 2},
		{GuestNumber: 1, Price: 150.0, Description: "Order 3", Checkin: time.Now(), Checkout: time.Now().Add(24 * time.Hour), CategoryType: types.ORIGINAL, RoomNumber: 103, UserID: 3},
		{GuestNumber: 3, Price: 250.99, Description: "Order 4", Checkin: time.Now(), Checkout: time.Now().Add(96 * time.Hour), CategoryType: types.VIP2, RoomNumber: 104, UserID: 4},
		{GuestNumber: 2, Price: 300.25, Description: "Order 5", Checkin: time.Now(), Checkout: time.Now().Add(120 * time.Hour), CategoryType: types.SEPARATELY, RoomNumber: 105, UserID: 5},
		{GuestNumber: 5, Price: 350.0, Description: "Order 6", Checkin: time.Now(), Checkout: time.Now().Add(144 * time.Hour), CategoryType: types.POPULAR, RoomNumber: 106, UserID: 6},
		{GuestNumber: 2, Price: 400.75, Description: "Order 7", Checkin: time.Now(), Checkout: time.Now().Add(168 * time.Hour), CategoryType: types.VIP1, RoomNumber: 107, UserID: 7},
		{GuestNumber: 1, Price: 450.5, Description: "Order 8", Checkin: time.Now(), Checkout: time.Now().Add(192 * time.Hour), CategoryType: types.ORIGINAL, RoomNumber: 108, UserID: 8},
		{GuestNumber: 3, Price: 500.25, Description: "Order 9", Checkin: time.Now(), Checkout: time.Now().Add(216 * time.Hour), CategoryType: types.POPULAR, RoomNumber: 109, UserID: 9},
		{GuestNumber: 4, Price: 550.0, Description: "Order 10", Checkin: time.Now(), Checkout: time.Now().Add(240 * time.Hour), CategoryType: types.VIP2, RoomNumber: 110, UserID: 10},
	}

	return data
}

func InitCategoriesDataDefault() []*enities.Category {

	data := []*enities.Category{
		{Name: "Original Suite", Description: "Standard room with basic amenities", ImageLink: "https://dq5r178u4t83b.cloudfront.net/wp-content/uploads/sites/125/2021/08/11060441/deluxe_harbour_web.jpg", Price: 100.0, AvailableRooms: 80, Type: types.ORIGINAL},
		{Name: "VIP Suite 1", Description: "VIP room with premium services", ImageLink: "https://assets.architecturaldigest.in/photos/65b2aecf269da4a0ee6c9b40/master/w_1600%2Cc_limit/atr.royalmansion-bedroom2-mr.jpg", Price: 200.0, AvailableRooms: 50, Type: types.VIP1},
		{Name: "VIP Suite 2", Description: "VIP room with extended premium services", ImageLink: "https://media.cnn.com/api/v1/images/stellar/prod/140127103345-peninsula-shanghai-deluxe-mock-up.jpg?q=w_2226,h_1449,x_0,y_0,c_fill", Price: 250.0, AvailableRooms: 25, Type: types.VIP2},
		{Name: "Popular Room", Description: "Popular choice among guests", ImageLink: "https://t3.ftcdn.net/jpg/06/19/00/08/360_F_619000872_AxiwLsfQqRHMkNxAbN4l5wg1MsPgBsmo.jpg", Price: 150.0, AvailableRooms: 100, Type: types.POPULAR},
		{Name: "Separately Room", Description: "Private room for individual guests", ImageLink: "https://t4.ftcdn.net/jpg/00/14/15/07/360_F_14150707_69W44NWjFECm7BXhnPCwX79XNTOSqEGm.jpg", Price: 180.0, AvailableRooms: 60, Type: types.SEPARATELY},
		{Name: "Deluxe Suite", Description: "Luxurious room with deluxe amenities", ImageLink: "https://www.cvent.com/sites/default/files/image/2021-10/hotel%20room%20with%20beachfront%20view.jpg", Price: 300.0, AvailableRooms: 70, Type: types.DELUXE},
	}

	return data
}

func InitRoomsDataDefault() []*enities.Room {

	data := []*enities.Room{
		{RoomNumber: 101, Status: types.RoomStatus(rand.Intn(3)), CategoryId: 1},
		{RoomNumber: 102, Status: types.RoomStatus(rand.Intn(3)), CategoryId: 1},
		{RoomNumber: 104, Status: types.RoomStatus(rand.Intn(3)), CategoryId: 2},
		{RoomNumber: 103, Status: types.RoomStatus(rand.Intn(3)), CategoryId: 3},
		{RoomNumber: 105, Status: types.RoomStatus(rand.Intn(3)), CategoryId: 4},
		{RoomNumber: 106, Status: types.RoomStatus(rand.Intn(3)), CategoryId: 6},
		{RoomNumber: 107, Status: types.RoomStatus(rand.Intn(3)), CategoryId: 5},
		{RoomNumber: 108, Status: types.RoomStatus(rand.Intn(3)), CategoryId: 4},
		{RoomNumber: 109, Status: types.RoomStatus(rand.Intn(3)), CategoryId: 2},
		{RoomNumber: 110, Status: types.RoomStatus(rand.Intn(3)), CategoryId: 3},
	}

	return data
}

func randomBirthday() *time.Time {
	start := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2005, 12, 31, 0, 0, 0, 0, time.UTC)
	delta := end.Sub(start)
	sec := rand.Int63n(int64(delta.Seconds()))
	randomTime := start.Add(time.Duration(sec) * time.Second)
	return &randomTime
}
