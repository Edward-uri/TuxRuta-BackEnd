package controller

import (
	"main/src/user/application"
	"main/src/user/domain/entities"

	"github.com/gin-gonic/gin"
)

type CreateUserHadler struct {
	CreateUserUseCase *application.CreateUserUseCase
}

func NewCreateUserHandler(createUserUseCase *application.CreateUserUseCase) *CreateUserHadler {
	return &CreateUserHadler{
		CreateUserUseCase: createUserUseCase,
	}
}

func (c *CreateUserHadler) HandleCreateUser(g *gin.Context) {
	var user entities.Usuario
	if err := g.ShouldBindJSON(&user); err != nil {
		g.JSON(400, gin.H{"error": "Invalid input data"})
		return
	}
	err := c.CreateUserUseCase.Execute(&user)
	if err != nil {
		g.JSON(500, gin.H{"error": err.Error()})
	}
	g.JSON(201, gin.H{"message": "User created successfully", "user": user})
}
