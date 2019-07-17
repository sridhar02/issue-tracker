describe("Repo Tests", () => {

	beforeEach(() => {
		cy.visit("/login")

		cy.get('[name="username"]').type("Sridhar1997")

		cy.get('[name="password"]').type("8082")

		cy.get('button').click();

	})

	it("Add and Remove Collaborator", () => {
		cy.visit('http://localhost:8000/Sridhar1997/funny/collaboration')

		cy.get('.search [name="user_name"]').type("DONE456")

		cy.get('.add-button').click();

		cy.url().should('eq', 'http://localhost:8000/Sridhar1997/funny/collaboration')


		cy.get('.close-button').click();			

		cy.url().should('eq', 'http://localhost:8000/Sridhar1997/funny/collaboration')
	})

	it("Non repo owner can not see collaboration page", () => {
		
		cy.visit('http://localhost:8000/DONE/test123/collaboration')

		cy.url().should('eq', 'http://localhost:8000/Sridhar1997')

	})

})