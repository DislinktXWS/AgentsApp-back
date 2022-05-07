package company_store

import "modules/dto"

func CommentMapper(commentReq *dto.RequestComment) Comment {
	comment := Comment{
		Content:   commentReq.Content,
		CompanyID: commentReq.CompanyID,
	}
	return comment
}
