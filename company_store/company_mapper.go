package company_store

import "modules/dto"

func companyMapper(companyReq *dto.RequestCompany) Company {
	company := Company{
		ID:                  companyReq.ID,
		CompanyCulture:      companyReq.CompanyCulture,
		Name:                companyReq.Name,
		YearOfEstablishment: companyReq.YearOfEstablishment,
		Address:             companyReq.Address,
		City:                companyReq.City,
		Country:             companyReq.Country,
		Phone:               companyReq.Phone,
		Industry:            companyReq.Industry,
		Description:         companyReq.Description,
		Email:               companyReq.Email,
		Website:             companyReq.Website,
		OwnerID:             companyReq.OwnerID,
	}
	return company
}
