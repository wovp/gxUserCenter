package gxmodule

import (
	"database/sql"
	"time"
)

// UserDetail结构体表示用户详细信息
type UserDetail struct {
	UserID      int
	FullName    string
	PhoneNumber string
	Address     string
	DateOfBirth time.Time
	Gender      string
	Occupation  string
	Avatar      string // 用户头像，这里假设存储头像的路径或URL
	Bio         string // 用户个人简介
	School      string // 用户学校
}

// Save将用户详细信息保存到数据库中
func (ud *UserDetail) Save(db *sql.DB) error {
	query := `
    INSERT INTO user_details (user_id, full_name, phone_number, address, date_of_birth, gender, occupation, avatar, bio, school)
    VALUES (?,?,?,?,?,?,?,?,?,?)
    `
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(ud.UserID, ud.FullName, ud.PhoneNumber, ud.Address, ud.DateOfBirth, ud.Gender, ud.Occupation, ud.Avatar, ud.Bio, ud.School)
	if err != nil {
		return err
	}

	_, err = result.LastInsertId()
	return err
}

// NewUserDetail创建一个新的UserDetail结构体实例
func NewUserDetail(fullName, phoneNumber, address, dateOfBirthStr, gender, occupation, avatar, bio, school string) (*UserDetail, error) {
	// 解析日期字符串为time.Time类型（假设日期格式为 "2000-01-01"）
	dateOfBirth, err := time.Parse("2006-01-02", dateOfBirthStr)
	if err != nil {
		return nil, err
	}

	return &UserDetail{
		FullName:    fullName,
		PhoneNumber: phoneNumber,
		Address:     address,
		DateOfBirth: dateOfBirth,
		Gender:      gender,
		Occupation:  occupation,
		Avatar:      avatar,
		Bio:         bio,
		School:      school,
	}, nil
}

// GetUserDetailByUserID根据用户ID从数据库中获取用户详细信息
func GetUserDetailByUserID(db *sql.DB, userID int) (*UserDetail, error) {
	query := `
    SELECT full_name, phone_number, address, date_of_birth, gender, occupation, avatar, bio, school
    FROM user_details
    WHERE user_id =?
    `
	row := db.QueryRow(query, userID)

	var userDetail UserDetail
	userDetail.UserID = userID

	err := row.Scan(&userDetail.FullName, &userDetail.PhoneNumber, &userDetail.Address, &userDetail.DateOfBirth, &userDetail.Gender, &userDetail.Occupation, &userDetail.Avatar, &userDetail.Bio, &userDetail.School)
	if err != nil {
		return nil, err
	}

	return &userDetail, nil
}

// UpdateUserDetail更新用户详细信息到数据库中
func (ud *UserDetail) UpdateUserDetail(db *sql.DB) error {
	query := `
	UPDATE user_details
	SET full_name =?, phone_number =?, address =?, date_of_birth =?, gender =?, occupation =?, avatar =?, bio =?, school =?
	WHERE user_id =?
	`
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(ud.FullName, ud.PhoneNumber, ud.Address, ud.DateOfBirth, ud.Gender, ud.Occupation, ud.Avatar, ud.Bio, ud.School, ud.UserID)
	if err != nil {
		return err
	}

	_, err = result.RowsAffected()
	if err != nil {
		return err
	}
	// 好奇怪，如果参数没变的话，他不更改
	//if rowsAffected == 0 {
	//	fmt.Println(result)
	//	return fmt.Errorf("没有找到对应的用户详细信息记录进行更新")
	//}

	return nil
}

// NewUserDetailRegister创建一个新的UserDetail结构体实例，可以选择使用默认值
func NewUserDetailRegister(fullName, phoneNumber, address, dateOfBirthStr, gender, occupation, avatar, bio, school string) (*UserDetail, error) {
	var defaultFullName, defaultPhoneNumber, defaultAddress, defaultGender, defaultOccupation, defaultAvatar, defaultBio, defaultSchool string
	var defaultDateOfBirth time.Time

	defaultFullName = "未填写姓名"
	defaultPhoneNumber = "未填写电话"
	defaultAddress = "未填写地址"
	defaultDateOfBirth = time.Time{}
	defaultGender = "未指定"
	defaultOccupation = "未填写职业"
	defaultAvatar = ""
	defaultBio = "暂无个人简介"
	defaultSchool = "未填写学校"

	// 解析日期字符串为time.Time类型（假设日期格式为 "2000-01-01"）
	var err error

	dateOfBirth := defaultDateOfBirth
	dateOfBirth, err = time.Parse("2000-01-01", dateOfBirthStr)
	if err != nil {
		return nil, err
	}

	return &UserDetail{
		FullName:    defaultFullName,
		PhoneNumber: defaultPhoneNumber,
		Address:     defaultAddress,
		DateOfBirth: dateOfBirth,
		Gender:      defaultGender,
		Occupation:  defaultOccupation,
		Avatar:      defaultAvatar,
		Bio:         defaultBio,
		School:      defaultSchool,
	}, nil
}
