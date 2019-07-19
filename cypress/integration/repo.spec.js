describe("Repo Tests", () => {

	beforeEach(() => {
		cy.visit("/login")
		cy.get('[name="username"]').type("Sridhar1997")
		cy.get('[name="password"]').type("8082")
		cy.get('button').click();
	})

	it("Add and Remove Collaborator", () => {
		cy.visit('/Sridhar1997/funny/collaboration')
		cy.get('.search [name="user_name"]').type("DONE456")
		cy.get('.add-button').click();
		cy.url().should('eq', 'http://localhost:8000/Sridhar1997/funny/collaboration')
		cy.get('.close-button').click();			
		cy.url().should('eq', 'http://localhost:8000/Sridhar1997/funny/collaboration')
	})

	it("Non repo owner can not see collaboration page", () => {
		cy.visit('/DONE/test123/collaboration')
		cy.url().should('eq', 'http://localhost:8000/Sridhar1997')
	})

	it("newrepo", () =>{
		const str = btoa(Math.random())	

		cy.visit("/Sridhar1997")
		cy.get('.newRepo-link').click()
		cy.get('[name="name"]').type("badshah" + str)
		cy.get('[name="description"]').type("songs in 2019")
		cy.get('[value="public"]').click()
		cy.get('.create-button').click()
		cy.get('.signout-button').click()
	})
})