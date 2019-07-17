describe("User related pages",  () => {

	it("sign in page ", () => {
		cy.visit("/login")

		cy.get('[name="username"]').type("Sridhar1997")

		cy.get('[name="password"]').type("8082")

		cy.get('button').click();

		cy.url().should('eq', 'http://localhost:8000/Sridhar1997')
	})

	it("sign in page (failure) ", () => {
		cy.visit("/login")

		cy.get('[name="username"]').type("Sridhar1997")

		cy.get('[name="password"]').type("BAD-PASSWORD")

		cy.get('button').click();

		cy.url().should('eq', 'http://localhost:8000/login')
	})

})