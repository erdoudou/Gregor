package controllers

type LogoutController struct {
	BaseController
}

func (c *LogoutController) Get() {
	c.DelSession("userLogin")
	c.Redirect("/login", 302)

}
