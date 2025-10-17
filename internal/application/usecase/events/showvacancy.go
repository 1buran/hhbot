package events

type ShowVacancy struct {
	Name, Salary, Skills, Desc, Employer, WorkFormat string
}

func NewShowVacancy(vacName, vacSalary, vacSkills, vacDesc, vacEmpl, vacWorkFormat string) ShowVacancy {
	return ShowVacancy{
		Name:       vacName,
		Salary:     vacSalary,
		Skills:     vacSkills,
		Desc:       vacDesc,
		Employer:   vacEmpl,
		WorkFormat: vacWorkFormat,
	}
}
