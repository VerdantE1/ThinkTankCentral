package response

import "ThinkTankCentral/model/database"

type FriendLinkInfo struct {
	List  []database.FriendLink `json:"list"`
	Total int64                 `json:"total"`
}
