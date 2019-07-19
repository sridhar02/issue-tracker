describe("newrepo pages", () => {


	it("newrepo", () =>{

	const str = btoa(Math.random())	

	cy.visit("/login")

	cy.get('[name="username"]').type("Sridhar1997")

	cy.get('[name="password"]').type("8082")

	cy.get('button').click()

	cy.visit("/Sridhar1997")

	cy.get('.newRepo-link').click()

	cy.get('[name="name"]').type("badshah" + str)

	cy.get('[name="description"]').type("songs in 2019")

	cy.get('[value="public"]').click()

	cy.get('.create-button').click()

	cy.get('.signout-button').click()

	})


})