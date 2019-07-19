describe("User related pages",  () => {

	it("sign out page ", () => {
		cy.visit("/login")

		cy.get('[name="username"]').type("Sridhar1997")

		cy.get('[name="password"]').type("8082")

		cy.get('button').click();

		cy.url().should('contains', '/Sridhar1997')

		cy.get('.signout-button').click()

		cy.url().should('contains', '/login')

	})


	it("sign up page", () => {

		const str = btoa(Math.random())

		cy.visit("/")

		cy.get('[name="name"]').type("dp" + str)
		
		cy.get('[name="username"]').type("dp21" +str)
		
		cy.get('[name="email"]').type("dp100" +str+ "@gmail.com")

		cy.get('[name="password"]').type("func2811")

		cy.get('.signup-button').click()

		cy.url().should('contains', '/login')
	})


})