describe("issue-tests repo owner", () =>{
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

	it("comment-close", () =>{
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

	it("issue-unpin", () =>{
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


describe("issue-tests collaborators",() =>{
	beforeEach(() => {
		cy.visit("/login")
		cy.get('[name="username"]').type("DONE456")
		cy.get('[name="password"]').type("x49b9LgbxdU7g2B")
		cy.get('button').click();
	})

	it("collaborator-close",() =>{
		cy.visit("/Sridhar1997/funny/issues/1")
		cy.get('[name="comment_and_close"]').click()	
	})

	it("collaborator-open",() =>{
		cy.visit("/Sridhar1997/funny/issues/1")
		cy.get('[name="comment_and_open"]').click()	
	})

	it("collaborator-pin", () =>{
		cy.visit("/Sridhar1997/funny/issues/1")
		cy.get('[name="pin-issue"]').click()
	})

	it("collaborator-unpin", () =>{
		cy.visit("/Sridhar1997/funny/issues/1")
		cy.get('[name="unpin-issue"]').click()
	})

	it("collaborator-lock",() => {
		cy.visit("/Sridhar1997/funny/issues/1")
		cy.get('[name="lock-button"]').click()
	})

	it("collaborator-unlock",() => {
		cy.visit("/Sridhar1997/funny/issues/1")
		cy.get('[name="unlock-button"]').click()
	})
})


describe("issue test of random",() =>{

	beforeEach(() => {
		cy.visit("/login")
		cy.get('[name="username"]').type("dp")
		cy.get('[name="password"]').type("func2811")
		cy.get('button').click()
	})

	it.only("new-issue",() =>{
		cy.visit("/Sridhar1997/funny/issues/1")
		cy.get('[href="/Sridhar1997/funny/issues/new"]').click()
		cy.get('[name="title"]').type("elen & mike")
		cy.get('[name="body"]').type("elven loves mike in stranger things")
		cy.get('.issue [type="submit"]').click()

		cy.visit("/Sridhar1997/funny/issues/5")
		cy.get('[name="body"]').type("elven loves mike")
		cy.get('.comment').click()

		cy.get('[name="comment_and_close"]').click()
		cy.get('.signout-button').click()


	})

})