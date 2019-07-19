describe("new issue and issue comment", () =>{


	beforeEach(() => {
		cy.visit("/login")
		cy.get('[name="username"]').type("Sridhar1997")
		cy.get('[name="password"]').type("8082")
		cy.get('button').click();
	})

	it("new issue" ,() =>{
		cy.visit("/Sridhar1997/Test/issues")
		cy.get('.new-issue').click()	
		cy.url().should('contains' ,"/Sridhar1997/Test/issues/new")

		cy.get('[name="title"]').type("elen & mike")
		cy.get('[name="body"]').type("elven loves mike in stranger things")
		cy.get('.issue [type="submit"]').click()

		cy.get('[name="comment_and_close"]').click()
	})

	it("new-comment", () =>{
		cy.visit("/Sridhar1997/Test/issues/7")
		cy.get('[name="body"]').type("elven loves mike")
		cy.get('.comment').click()
	})

	it.only("comment-close", () =>{
		cy.visit("/Sridhar1997/Test/issues/7")
		cy.get('[name="comment_and_close"]').click()
	})

	it("comment-open", () =>{
		cy.visit("/Sridhar1997/Test/issues/7")
		cy.get('[name="comment_and_open"]').click()
	})


	it("issue-pin", () =>{
		cy.visit("/Sridhar1997/Test/issues/7")
		cy.get('[name="pin-issue"]').click()
	})

	it("issue-pin", () =>{
		cy.visit("/Sridhar1997/Test/issues/7")
		cy.get('[name="unpin-issue"]').click()
	})

	it("issue-lock",() => {
		cy.visit("/Sridhar1997/Test/issues/7")
		cy.get('[name="lock-button"]').click()
	})

	it("issue-unlock",() => {
		cy.visit("/Sridhar1997/Test/issues/7")
		cy.get('[name="unlock-button"]').click()
	})





})