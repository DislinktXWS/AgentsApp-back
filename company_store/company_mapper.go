package company_store

import (
	"modules/dto"
)

func CompanyMapper(companyReq *dto.RequestCompany) Company {
	company := Company{
		CompanyCulture:      companyReq.CompanyCulture,
		Name:                companyReq.Name,
		YearOfEstablishment: companyReq.YearOfEstablishment,
		Address:             companyReq.Address,
		Phone:               companyReq.Phone,
		Industry:            companyReq.Industry,
		Description:         companyReq.Description,
		Email:               companyReq.Email,
		Website:             companyReq.Website,
		OwnerID:             companyReq.OwnerID,
		Accepted:            false,
		Checked:             false,
	}
	return company
}
