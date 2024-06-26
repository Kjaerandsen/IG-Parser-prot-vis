package parser

import (
	"fmt"
	"regexp"
	"testing"
)

/*
Testing basic component syntax, including suffices and annotations,
with particular focus on special characters.
*/
func TestSingleComponentSyntaxSpecialCharacters(t *testing.T) {

	// Implies use of special characters
	r, err := regexp.Compile(FULL_COMPONENT_SYNTAX)
	if err != nil {
		t.Fatal("Error during compilation:", err.Error())
	}

	text := "(Bklsdjgl#k{sgv sk}lvjds) dalsjglks() Bdir,p1[ruler=gove§rnor](jglkdsjgsiovs) " +
		"Cac[left=right[anotherLeft,an@otherRight],right=[left,right], key=values]{A(actor) I(aim)} " +
		"P,p343(ano@ther comp€onent values#$) " +
		"E1{ A(acto¤§r two) I1(aim1) }" +
		" text outside"

	res := r.FindAllString(text, -1)

	fmt.Println("Matching component structure (primitive and nested)")
	fmt.Println(res)
	fmt.Println("Count:", len(res))

	firstElem := "Bdir,p1[ruler=gove§rnor](jglkdsjgsiovs)"

	if res[0] != firstElem {
		t.Fatal("Wrong element matched. Should be", firstElem, ", but is "+res[0])
	}

	secondElem := "Cac[left=right[anotherLeft,an@otherRight],right=[left,right], key=values]{A(actor) I(aim)}"

	if res[1] != secondElem {
		t.Fatal("Wrong element matched. Should be", secondElem, ", but is "+res[1])
	}

	thirdElem := "P,p343(ano@ther comp€onent values#$)"

	if res[2] != thirdElem {
		t.Fatal("Wrong element matched. Should be", thirdElem, ", but is "+res[2])
	}

	fourthElem := "E1{ A(acto¤§r two) I1(aim1) }"

	if res[3] != fourthElem {
		t.Fatal("Wrong element matched. Should be", fourthElem, ", but is "+res[3])
	}

	if len(res) != 4 {
		t.Fatal("Wrong number of matched elements. Should be 4, but is", len(res))
	}
}

/*
Testing basic component syntax, including suffices and annotations
*/
func TestSingleComponentSyntax(t *testing.T) {

	r, err := regexp.Compile(FULL_COMPONENT_SYNTAX)
	if err != nil {
		t.Fatal("Error during compilation:", err.Error())
	}

	text := "(Bklsdjgl#k{sgv sk}lvjds) dalsjglks() Bdir,p1[ruler=governor](jglkdsjgsiovs) " +
		"Cac[left=right[anotherLeft,anotherRight],right=[left,right], key=values]{A(actor) I(aim)} " +
		"P,p343(another component values#$) " +
		"E1{ A(actor two) I1(aim1) }" +
		" text outside"

	res := r.FindAllString(text, -1)

	fmt.Println("Matching component structure (primitive and nested)")
	fmt.Println(res)
	fmt.Println("Count:", len(res))

	firstElem := "Bdir,p1[ruler=governor](jglkdsjgsiovs)"

	if res[0] != firstElem {
		t.Fatal("Wrong element matched. Should be", firstElem, ", but is "+res[0])
	}

	secondElem := "Cac[left=right[anotherLeft,anotherRight],right=[left,right], key=values]{A(actor) I(aim)}"

	if res[1] != secondElem {
		t.Fatal("Wrong element matched. Should be", secondElem, ", but is "+res[1])
	}

	thirdElem := "P,p343(another component values#$)"

	if res[2] != thirdElem {
		t.Fatal("Wrong element matched. Should be", thirdElem, ", but is "+res[2])
	}

	fourthElem := "E1{ A(actor two) I1(aim1) }"

	if res[3] != fourthElem {
		t.Fatal("Wrong element matched. Should be", fourthElem, ", but is "+res[3])
	}

	if len(res) != 4 {
		t.Fatal("Wrong number of matched elements. Should be 4, but is", len(res))
	}
}

/*
Tests for combinations within text. Note that is does not test for terminated statement combinations. That is tested in statement parsing tests.
*/
func TestComponentCombinations(t *testing.T) {

	// Note: Only used in testing; in production NESTED_COMBINATIONS_TERMINATED is used
	r, err := regexp.Compile(NESTED_COMBINATIONS)
	if err != nil {
		t.Fatal("Error during compilation:", err.Error())
	}

	text := "(Aklsdjgl#k{sgv sk}lvjds) {[]jdskgl ds()} Bdir,p1[ruler=governor](jglkdsjgsiovs) Cac[left=right[anotherLeft,anotherRight],right=[left,right], key=values]{A(actor) I(aim)}" +
		"Cac{A(dlkgjsg) I[dgisg](kjsdglkds) [AND] (Bdir{djglksjdgkd} Cex(A(sdlgjlskd)) [XOR] A(dsgjslkj) E(gklsjgls))}" +
		"Cac{Cac{ A(actor) I(fjhgjh) Bdir(rtyui)} [XOR] Cac{A(ertyui) I(dfghj)}}" +
		"Cac{Cac{ A(as(dslks)a) I(adgklsjlg)} [XOR] Cac(asas) [AND] Cac12[kgkg]{lkdjgdls} [OR] A(dslgkjds)}" +
		"Cac{Cac(andsdjsglk) [AND] A(sdjlgsl) Bdir(jslkgsjlkgds)}" +
		"Cac{Cac(andsdjsglk) [AND] ( A(sdjlgsl) [XOR] (A(sdoidjs) [OR] A(sdjglksj)))}" +
		"((dglkdsjg [AND] jdlgksjlkgd))"

	res := r.FindAllString(text, -1)

	fmt.Println("Refined matching combinations")
	fmt.Println(res)
	fmt.Println("Count:", len(res))

	firstElem := "Cac{A(dlkgjsg) I[dgisg](kjsdglkds) [AND] (Bdir{djglksjdgkd} Cex(A(sdlgjlskd)) [XOR] A(dsgjslkj) E(gklsjgls))}"

	if res[0] != firstElem {
		t.Fatal("Wrong element matched. Should be", firstElem, ", but is "+res[0])
	}

	secondElem := "Cac{Cac{ A(actor) I(fjhgjh) Bdir(rtyui)} [XOR] Cac{A(ertyui) I(dfghj)}}"

	if res[1] != secondElem {
		t.Fatal("Wrong element matched. Should be", secondElem, ", but is "+res[1])
	}

	thirdElem := "Cac{Cac{ A(as(dslks)a) I(adgklsjlg)} [XOR] Cac(asas) [AND] Cac12[kgkg]{lkdjgdls} [OR] A(dslgkjds)}"

	if res[2] != thirdElem {
		t.Fatal("Wrong element matched. Should be", thirdElem, ", but is "+res[2])
	}

	fourthElem := "Cac{Cac(andsdjsglk) [AND] A(sdjlgsl) Bdir(jslkgsjlkgds)}"

	if res[3] != fourthElem {
		t.Fatal("Wrong element matched. Should be", fourthElem, ", but is "+res[3])
	}

	fifthElem := "Cac{Cac(andsdjsglk) [AND] ( A(sdjlgsl) [XOR] (A(sdoidjs) [OR] A(sdjglksj)))}"

	if res[4] != fifthElem {
		t.Fatal("Wrong element matched. Should be", fifthElem, ", but is "+res[4])
	}

	if len(res) != 5 {
		t.Fatal("Wrong number of matched elements. Should be 5, but is", len(res))
	}

}

/*
Tests for component pair combinations within text.
*/
func TestComponentPairCombinations(t *testing.T) {

	// Note: Only used in testing; in production NESTED_COMBINATIONS_TERMINATED is used
	r, err := regexp.Compile(COMPONENT_PAIR_COMBINATIONS)
	if err != nil {
		t.Fatal("Error during compilation:", err.Error())
	}

	text := "(Aklsdjgl#k{sgv sk}lvjds) {[]jdskgl ds()} Bdir,p1[ruler=governor](jglkdsjgsiovs) Cac[left=right[anotherLeft,anotherRight],right=[left,right], key=values]{A(actor) I(aim)}" +
		"{A(dlkgjsg) I[dgisg](kjsdglkds) [AND] (Bdir{djglksjdgkd} Cex(A(sdlgjlskd)) [XOR] A(dsgjslkj) E(gklsjgls))}" +
		"Cac{Cac{ A(actor) I(fjhgjh) Bdir(rtyui)} [XOR] Cac{A(ertyui) I(dfghj)}}" +
		"[dfghjk]{Cac{ A(as(dslks)a) I(adgklsjlg)} [XOR] Cac(asas) [AND] Cac12[kgkg]{lkdjgdls} [OR] A(dslgkjds)}" +
		"{Cac(andsdjsglk) [AND] A(sdjlgsl) Bdir(jslkgsjlkgds)}" +
		"{Cac(andsdjsglk) [AND] ( A(sdjlgsl) [XOR] (A(sdoidjs) [OR] A(sdjglksj)))}" +
		" {Cac{A(sdlkgjsdlk) Bdir{A(aljdgs) I(kdsjglkj) Bdir(dkslgj)}} [XOR] Cac{A(skdfjcs) Bdir{A(dlksgjie) I(dsklgjiv) Bdir(lkdsjgei)}}} " +
		"((dglkdsjg [AND] jdlgksjlkgd))"

	res := r.FindAllString(text, -1)

	fmt.Println("Refined matching combinations")
	fmt.Println(res)
	fmt.Println("Count:", len(res))

	firstElem := "{A(dlkgjsg) I[dgisg](kjsdglkds) [AND] (Bdir{djglksjdgkd} Cex(A(sdlgjlskd)) [XOR] A(dsgjslkj) E(gklsjgls))}"

	if res[0] != firstElem {
		t.Fatal("Wrong element matched. Should be", firstElem, ", but is "+res[0])
	}

	secondElem := "{Cac{ A(actor) I(fjhgjh) Bdir(rtyui)} [XOR] Cac{A(ertyui) I(dfghj)}}"

	if res[1] != secondElem {
		t.Fatal("Wrong element matched. Should be", secondElem, ", but is "+res[1])
	}

	thirdElem := "[dfghjk]{Cac{ A(as(dslks)a) I(adgklsjlg)} [XOR] Cac(asas) [AND] Cac12[kgkg]{lkdjgdls} [OR] A(dslgkjds)}"

	if res[2] != thirdElem {
		t.Fatal("Wrong element matched. Should be", thirdElem, ", but is "+res[2])
	}

	fourthElem := "{Cac(andsdjsglk) [AND] A(sdjlgsl) Bdir(jslkgsjlkgds)}"

	if res[3] != fourthElem {
		t.Fatal("Wrong element matched. Should be", fourthElem, ", but is "+res[3])
	}

	fifthElem := "{Cac(andsdjsglk) [AND] ( A(sdjlgsl) [XOR] (A(sdoidjs) [OR] A(sdjglksj)))}"

	if res[4] != fifthElem {
		t.Fatal("Wrong element matched. Should be", fifthElem, ", but is "+res[4])
	}

	sixthElem := "{Cac{A(sdlkgjsdlk) Bdir{A(aljdgs) I(kdsjglkj) Bdir(dkslgj)}} [XOR] Cac{A(skdfjcs) Bdir{A(dlksgjie) I(dsklgjiv) Bdir(lkdsjgei)}}}"

	if res[5] != sixthElem {
		t.Fatal("Wrong element matched. Should be", sixthElem, ", but is "+res[5])
	}

	if len(res) != 6 {
		t.Fatal("Wrong number of matched elements. Should be 6, but is", len(res))
	}
}

/*
Tests complex statement combinations that reflect nesting characteristics.
*/
func TestComplexStatementCombinations(t *testing.T) {

	text := " {  Cac{A(actor1) I(aim1) Bdir(object1)}   [AND]   Cac{A(actor2)  I(aim2) Bdir(object2) }   } "
	text += "{{Cac{ A(actor1) I(aim1) Bdir(object1) }   [XOR]  Cac{ fgfd A(actor1a) fdhdf I(aim1a) Bdir(object1a)}} dfsjfdsl [AND] lkdsjflksj {Cac{A(actor2) I(aim2) Bdir(object2)} [OR] Cac{A(actor3) I(aim3) Bdir(object3)}}}"
	text += " A(dfkflslkjfs) Cac(dlsgjslkdj) " // should not be found
	text += "{{{Cac{ A(actor1) I(aim1) Bdir(object1) } [XOR] Cac{ A(actor1) I(aim1) Bdir(object1) }}   [XOR]  Cac{ fgfd A(actor1a) fdhdf I(aim1a) Bdir(object1a)}} dfsjfdsl [AND] lkdsjflksj {{Cac{A(actor2) I(aim2) Bdir(object2)} dgjsksldgj[XOR] Cac{A(actor2) I(aim2) Bdir(object2)}} [OR] Cac{A(actor3) I(aim3) Bdir(object3)}}}"
	text += "{Cac{A(actor1) I(aim1) Bdir{A(actor2) I(aim2) Cac(condition2)}} [OR] Cac{A(actor3) I(aim3) Bdir(object3)}}"
	text += "{blabla ,.Cac{A(actor1) I(aim1) Bdir{A(actor2) I(aim2) Cac(condition2)}}, [OR]Cac{A(actor3) I(aim3) Bdir(object3)}more; text}"

	r, err := regexp.Compile(BRACED_6TH_ORDER_COMBINATIONS)
	if err != nil {
		t.Fatal("Error during compilation:", err.Error())
	}

	res := r.FindAllString(text, -1)

	if len(res) != 5 {
		outString := ""
		for i, v := range res {
			item := fmt.Sprint("Item ", i, ": Value: ", v)
			outString += item + "\n"
		}
		t.Fatal("Number of statements is not correct. Should be 5, but is", len(res), "\nStatements:\n"+outString)
	}

	firstElem := "{  Cac{A(actor1) I(aim1) Bdir(object1)}   [AND]   Cac{A(actor2)  I(aim2) Bdir(object2) }   }"

	if res[0] != firstElem {
		t.Fatal("Element incorrect. It should read '"+firstElem+"', but is", res[0])
	}

	secondElem := "{{Cac{ A(actor1) I(aim1) Bdir(object1) }   [XOR]  Cac{ fgfd A(actor1a) fdhdf I(aim1a) Bdir(object1a)}} dfsjfdsl [AND] lkdsjflksj {Cac{A(actor2) I(aim2) Bdir(object2)} [OR] Cac{A(actor3) I(aim3) Bdir(object3)}}}"

	if res[1] != secondElem {
		t.Fatal("Element incorrect. It should read '"+secondElem+"', but is", res[1])
	}

	thirdElem := "{{{Cac{ A(actor1) I(aim1) Bdir(object1) } [XOR] Cac{ A(actor1) I(aim1) Bdir(object1) }}   [XOR]  Cac{ fgfd A(actor1a) fdhdf I(aim1a) Bdir(object1a)}} dfsjfdsl [AND] lkdsjflksj {{Cac{A(actor2) I(aim2) Bdir(object2)} dgjsksldgj[XOR] Cac{A(actor2) I(aim2) Bdir(object2)}} [OR] Cac{A(actor3) I(aim3) Bdir(object3)}}}"

	if res[2] != thirdElem {
		t.Fatal("Element incorrect. It should read '"+thirdElem+"', but is", res[2])
	}

	fourthElem := "{Cac{A(actor1) I(aim1) Bdir{A(actor2) I(aim2) Cac(condition2)}} [OR] Cac{A(actor3) I(aim3) Bdir(object3)}}"

	if res[3] != fourthElem {
		t.Fatal("Element incorrect. It should read '"+fourthElem+"', but is", res[3])
	}

	// Tests for tolerance toward comma following logical operator, missing separating space, and excessive terms
	fifthElem := "{blabla ,.Cac{A(actor1) I(aim1) Bdir{A(actor2) I(aim2) Cac(condition2)}}, [OR]Cac{A(actor3) I(aim3) Bdir(object3)}more; text}"

	if res[4] != fifthElem {
		t.Fatal("Element incorrect. It should read '"+fifthElem+"', but is", res[4])
	}
}

/*
Tests complex statement combinations that reflect nesting characteristics, but also embedded
diacritics in content as well as annotations.
*/
func TestComplexStatementCombinationsWithDiacriticsAndAnnotations(t *testing.T) {

	text := " {  Cac{A[äsdlj](actor1) I[sdgjöñ](aim1) Bdir(objäect1)}   [AND]  ï Cac{A(acřtor2)  I(aim2) Bdir(object2) }   } "
	text += "{{Cac{ A(actōor1) I[dlksjëjvcls](aimë1) Bdir(object1) }   [XOR]  Cac{ éfgfd A[djösgdē](acåtor1a) fdhdf I(aim1a) Bdir(object1a)}} dfsjŏfdsl [AND] lkdsjflksj {Cac{A(actoûr2) I(aiŏm2) Bdir(object2)} [OR] Cac{A(actor3) I(aim3) Bdir(object3)}}}"
	text += " A(dfkflslküjfs) Cac[slkÅpså](dlsgjslkdj) ô ö Ü ä" // should not be found
	text += "{{{Cac{ A(acÜtor1) I(aim1) Bdir(objecät1) } [XOR] Cac[ÅlidsLøØ]{ A(actor1) I(aim1) Bdir(object1) }}   [XOR]  Cac{ fgfd A[randsléø](actor1a) fdhdf I(aéim1a) Bdir(object1a)}} dfsjfdsl [AND] lkdsjflksj {{Cac{A(actor2) I(aim2) Bdir(object2)} dgjsksldgj[XOR] Cac{A(actor2) I(aim2) Bdir(object2)}} [OR] Cac{A(actor3) I(aim3) Bdir(object3)}}}"
	text += "{Cac{A(actoşr1) I(aïim1) Bdir[ddlgžöjř]{A(actor2) I(aim2) Cac(condñition2)}} [OR] Cac{A(actáor3) I(aim3) Bdir(object3)}}"
	text += "{blabla ,.Cac{A(actôor1) I(aim1) Bdir{A(actor2) I(aim2) Cac(coçndition2)}}, [OR]Cac{A(aēctor3) I[dsklfjöçjsdlfö](aim3) Bdir(object3)}more; ïtext}"

	r, err := regexp.Compile(BRACED_6TH_ORDER_COMBINATIONS)
	if err != nil {
		t.Fatal("Error during compilation:", err.Error())
	}

	res := r.FindAllString(text, -1)

	if len(res) != 5 {
		outString := ""
		for i, v := range res {
			item := fmt.Sprint("Item ", i, ": Value: ", v)
			outString += item + "\n"
		}
		t.Fatal("Number of statements is not correct. Should be 5, but is", len(res), "\nStatements:\n"+outString)
	}

	firstElem := "{  Cac{A[äsdlj](actor1) I[sdgjöñ](aim1) Bdir(objäect1)}   [AND]  ï Cac{A(acřtor2)  I(aim2) Bdir(object2) }   }"

	if res[0] != firstElem {
		t.Fatal("Element incorrect. It should read '"+firstElem+"', but is", res[0])
	}

	secondElem := "{{Cac{ A(actōor1) I[dlksjëjvcls](aimë1) Bdir(object1) }   [XOR]  Cac{ éfgfd A[djösgdē](acåtor1a) fdhdf I(aim1a) Bdir(object1a)}} dfsjŏfdsl [AND] lkdsjflksj {Cac{A(actoûr2) I(aiŏm2) Bdir(object2)} [OR] Cac{A(actor3) I(aim3) Bdir(object3)}}}"

	if res[1] != secondElem {
		t.Fatal("Element incorrect. It should read '"+secondElem+"', but is", res[1])
	}

	thirdElem := "{{{Cac{ A(acÜtor1) I(aim1) Bdir(objecät1) } [XOR] Cac[ÅlidsLøØ]{ A(actor1) I(aim1) Bdir(object1) }}   [XOR]  Cac{ fgfd A[randsléø](actor1a) fdhdf I(aéim1a) Bdir(object1a)}} dfsjfdsl [AND] lkdsjflksj {{Cac{A(actor2) I(aim2) Bdir(object2)} dgjsksldgj[XOR] Cac{A(actor2) I(aim2) Bdir(object2)}} [OR] Cac{A(actor3) I(aim3) Bdir(object3)}}}"

	if res[2] != thirdElem {
		t.Fatal("Element incorrect. It should read '"+thirdElem+"', but is", res[2])
	}

	fourthElem := "{Cac{A(actoşr1) I(aïim1) Bdir[ddlgžöjř]{A(actor2) I(aim2) Cac(condñition2)}} [OR] Cac{A(actáor3) I(aim3) Bdir(object3)}}"

	if res[3] != fourthElem {
		t.Fatal("Element incorrect. It should read '"+fourthElem+"', but is", res[3])
	}

	// Tests for tolerance toward comma following logical operator, missing separating space, and excessive terms
	fifthElem := "{blabla ,.Cac{A(actôor1) I(aim1) Bdir{A(actor2) I(aim2) Cac(coçndition2)}}, [OR]Cac{A(aēctor3) I[dsklfjöçjsdlfö](aim3) Bdir(object3)}more; ïtext}"

	if res[4] != fifthElem {
		t.Fatal("Element incorrect. It should read '"+fifthElem+"', but is", res[4])
	}
}
