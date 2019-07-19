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


})