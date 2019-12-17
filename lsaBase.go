// script to compute  SVD of a Term Document Matrix created from a text corpus

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/gonum/matrix/mat64"
	//"github.com/gonum/matrix"
)

func main() {

	//commWords is from the corpus of docs with potential keywords deleted leaving only
	//common words that are unlikely to be used as terms in LSA.
	textCommon := rCommon.Replace(commWords) //remove punctuation, replace by spaces

	commList := DupCheck(strings.Fields(strings.ToLower(textCommon)))
	sort.Strings(commList)
	//Now to strip away these words from the documents
	stripList := make([]string, 0)
	for i := 0; i < len(commList); i++ {
		stripList = append(stripList, "&"+commList[i]+"&")
		stripList = append(stripList, " ")
	}

	//list of documents as array/slice of original string constants
	docList := [...]string{doc1, doc2, doc3, doc4}

	//extract potential tokens
	docs := make([]string, len(docList))
	//combined tokens (combine all tokens together)
	combinedTerms := make([]string, 0)
	for i := 0; i < len(docList); i++ {
		tempToks, _ := getTokens(docList[i], stripList)
		docs[i] = " " + strings.Join(tempToks, " ") + " " // add a space at the
		//start and end for counting
		//otherwise, last term is not counted with " "+term+" "
		combinedTerms = append(combinedTerms, tempToks...)
	}
	//Remove duplicate terms
	terms := DupCheck(combinedTerms) //
	//sort the list
	sort.Strings(terms)

	//Now to go about creating the document terms matrix
	//Use of append because I do not know how to initialize multidimensional slices when the length is unknown.
	zeros := []int{0}
	for i := 0; i < len(docs)-1; i++ { //len-1 because first element is already initialized
		zeros = append(zeros, []int{0}...)
	}

	//document term matrix (initialization)
	dtm := [][]int{zeros}
	//create an empty term document matrix
	for i := 0; i < len(terms)-1; i++ { // len-1 because first row is already initialized
		dtm = append(dtm, [][]int{zeros}...)
	}
	//now assign elements for the term-document matrix
	fmt.Printf("\nTerm document matrix with %d terms.\n", len(terms))
	for i := 0; i < len(terms); i++ {
		fmt.Printf("\n%d %v: ", i+1, terms[i])
		for j := 0; j < len(docs); j++ {
			//dtm[i][j] = strings.Count(docs[j],terms[i])
			dtm[i][j] = strings.Count(docs[j], " "+terms[i]+" ")
			fmt.Printf("%d ", dtm[i][j])
		}
	}

	//DTM with restricted keywords
	keys := strings.Fields(keyWords)
	dtm2 := [][]int{zeros}
	//create an empty term document matrix
	for i := 0; i < len(keys)-1; i++ { // len-1 because first row is already initialized
		dtm2 = append(dtm2, [][]int{zeros}...)
	}
	//create a new float matrix
	a := mat64.NewDense(len(keys), len(docs), nil)

	//assign elements for the term-document matrix
	fmt.Printf("\n\nKeyword-based Term document matrix with %d terms.\n", len(keys))
	for i := 0; i < len(keys); i++ {
		fmt.Printf("\n%d %v: ", i+1, keys[i])
		for j := 0; j < len(docs); j++ {
			//dtm[i][j] = strings.Count(docs[j],terms[i])
			dtm2[i][j] = strings.Count(docs[j], keys[i])
			a.Set(i, j, float64(dtm2[i][j]))
			fmt.Printf("%d ", dtm2[i][j])
		}
	}

	fmt.Printf("\n")

	// Create a matrix formatting value with a prefix ...
	fa := mat64.Formatted(a, mat64.Prefix("    "))

	// and then print with and without zero value elements.
	fmt.Printf("with all values:\na = %v\n\n", fa)
	//fmt.Printf("with only non-zero values:\na = % v\n\n", fa)

	// Compute the SVD factorization.
	var svd mat64.SVD
	if ok := svd.Factorize(a, 2); !ok { // 2 implies SVDThin
		fmt.Println("Problem with SVD")
	}
	SVs := svd.Values(nil)
	numSVs := len(svd.Values(nil))
	fmt.Printf("\nThere are %d singular values %v\n", numSVs, SVs)
	//dimensions
	r, c := a.Dims()
	//U and V-transpose from SVD
	U := mat64.NewDense(r, numSVs, nil)
	U.UFromSVD(&svd)                                                      //This is U
	fmt.Printf("\nU = %0.4v\n", mat64.Formatted(U, mat64.Prefix("    "))) //4 decimal places
	V := mat64.NewDense(c, numSVs, nil)
	V.VFromSVD(&svd) //this is V^T
	fmt.Printf("\nV = %0.4v\n", mat64.Formatted(V, mat64.Prefix("    ")))
	//fmt.Printf("\nThe (1,2) element of V^T is %0.6v\n",V.At(1,2))

}

//Tokenizer function remove common words, print out sorted tokens
func getTokens(input string, replacer []string) ([]string, error) {
	toks := make([]string, 0)
	text := strings.ToLower(string(input))
	text = rCommon.Replace(text)
	text = r1.Replace(text)
	//define a new replacer
	var rnew = strings.NewReplacer(replacer...) //didn't work. Needs to be interface or
	//when newplreplacer was not working, back up
	//for i:=0;i<len(replacer);i++{
	//	text = strings.Replace(text,replacer[i], " ", -1)
	//}
	text = rnew.Replace(text)
	//text = r2.Replace(text)

	text = r3.Replace(text)

	//NewReader returns a new Reader reading from s. It is similar to
	//bytes.NewBufferString but more efficient and read-only.
	scanner := bufio.NewScanner(strings.NewReader(text))
	// Set the split function for the scanning operation.
	//If called, it must be called before Scan. The default split function is ScanLines.
	scanner.Split(bufio.ScanWords) //the split function here is ScanWords
	// Count the words.
	//count := 0
	//toks := make([]string,0) //slice of tokens

	//func (s *Scanner) Scan() bool
	//Scan advances the Scanner to the next token, which will then be available through
	//the Bytes or Text method. It returns false when the scan stops, either by reaching
	//the end of the input or an error. After Scan returns false, the Err method will
	//return any error that occurred during scanning, except that if it was io.EOF, Err
	//will return nil. Scan panics if the split function returns 100 empty tokens without
	//advancing the input. This is a common error mode for scanners.
	for scanner.Scan() { //calls the scanner function
		toks = append(toks, scanner.Text())
		//fmt.Printf("%v\n",toks[count])
		//count++
	}
	var err error
	if err = scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err)
	}
	//if sorting makes sense
	//sort.Strings(toks)
	return toks, err
}

//check and remove exact duplicates
func DupCheck(s []string) (sOut []string) {
	sStr := " " + strings.Join(s, " ") + " " //join with spaces, add initial and final space
	numStr := len(s)
	for i := 0; i < numStr; i++ {
		dup := strings.Count(sStr, " "+s[i]+" ") //look for exact match not submatch
		if dup > 1 {
			sStr = strings.Replace(sStr, " "+s[i]+" ", " ", dup-1) //replace all but one.
		}
	}
	sOut = strings.Fields(sStr) //convert back to []string
	return
}

// Create "Replacer"s

// to replace punctuation etc by spaces.
var rCommon = strings.NewReplacer(", ", " ", `"`, " ", `. `, " ",
	"-\n", "", "-\r", "", //line break hyphens
	"-", " ", //this is a regular hyphen
	"–", " ", //this is a PDF ligature hyphen (e.g. when using pdf2txt)
	":", " ", ";", " ",
	"(", " ", ") ", " ", "), ", " ", "). ", " ", ".\n", " ",
	"\n", " ", "\t", " ", "\r", " ")

//replace all spaces by &&
var r1 = strings.NewReplacer(" ", "&&")

// r2 is created programmatically.

var r3 = strings.NewReplacer("&", " ")

// THE CORPUS: In actual practice, these would be separate files extracted from PDFs (for journal articles) or parsed HTML

const doc1 string = ` Patients after solid organ transplantation (SOT) carry a
	substantially	increased risk to develop malignant lymphomas. This is in part due to the 
 immunosuppression required to maintain the function of the organ graft. 
 Depending on the transplanted organ, up to 15% of pediatric transplant recipients acquire 
 posttransplant lymphoproliferative disease (PTLD), and eventually 20% of those succumb to 
 the disease. Early diagnosis of PTLD is often hampered by the unspecific symptoms and the 
 difficult differential diagnosis, which includes atypical infections as well as graft  
 rejection. Treatment of PTLD is limited by the high vulnerability towards antineoplastic 
 chemotherapy in transplanted children. However, new treatment strategies and especially 
 the introduction of the monoclonal anti-CD20 antibody rituximab have dramatically 
  improved 
  outcomes of PTLD. This review discusses risk factors for the development of PTLD in 
  children, summarizes current approaches to therapy, and gives an outlook on developing 
  new treatment modalities like targeted therapy with virus-specific T cells. 
  Finally, monitoring strategies are evaluated.  `

const doc2 string = ` PTLD induced by Epstein Barr Virus (EBV) is a major cause of 
 morbidity and mortality in solid organ recipients. Despite a growing understanding of the
  pathogenesis of EBV infection and EBV-associated diseases in transplant recipients there 
  remains uncertainty regarding the best clinical management of these patients.
 These guidelines are offered to improve care by establishing consistent evidence-based 
 care.  Where evidence in the pediatric population is lacking, literature from adult 
 patients was reviewed and used where applicable and/or local consensus was used to form
 the recommendations. These guidelines do not address EBV negative PTLD. `

const doc3 string = ` A joint working group established by the Haemato-oncology subgroup
  of the British Committee for Standards in Haematology (BCSH) and the British 
  Transplantation Society (BTS) has reviewed the available literature and made 
  recommendations for the diagnosis and management of post-transplant lymphoproliferative 
  disorder in adult recipients of solid organ transplants. This review details the 
  therapeutic options recommended including reduction in immunosuppression (RIS), 
  transplant organ resection, radiotherapy and chemotherapy. Effective therapy should be 
  instituted before progressive disease results in declining performance status and 
  multiorgan dysfunction. The goal of treatment should be a durable complete remission 
  with retention of transplanted organ function with minimal toxicity. `

const doc4 string = ` Post-transplant lymphoproliferative disorder (PTLD) represents a 
 spectrum of Epstein–Barr virus-related (EBV) clinical diseases, from a benign 
 mononucleosis-like illness to a fulminant non-Hodgkin’s lymphoma. In the setting of 
 hematopoietic stem cell transplantation, PTLD is an often-fatal complication occurring 
 relatively early after transplant. Risk factors for the development of PTLD are well 
 established, and include HLA-mismatching, T-cell depletion, and the use of 
 antilymphocyte antibodies as conditioning or treatment of graft-versus- host disease. 
 Early recognition of PTLD is particularly important in the SCT setting, because PTLD in 
 these patients tends to be rapidly progressive. Familiarity with the clinical features of 
 PTLD and a heightened level of suspicion are critical for making the diagnosis. 
 Surveillance techniques with EBV antibody titers and/or polymerase chain reaction (PCR) 
 may have a role in some high-risk settings. Immune-based therapies such as monoclonal 
 anti-B-cell antibodies, interferon*A, and EBV-specific donor T cells, either as treatment
  for PTLD or as prophylaxis in high-risk patients, represent promising new directions in 
 the treatment of this disease. `

const commWords string = ` Patients after  carry a 
substantially increased risk to develop   This is in part due to the 
 required to maintain the function of the  
Depending on the  up to of acquire 
 and eventually of those succumb to 
the disease Early diagnosis of  is often hampered by the unspecific  and the 
difficult differential  which includes atypical as well as   
Treatment of  is limited by the high vulnerability towards  However new treatment strategies and especially 
the introduction of the have dramatically improved
 outcomes of  This review discusses risk factors for the development of  in 
 summarizes current approaches to therapy and gives an outlook on developing 
 new treatment modalities like targeted therapy with 
 Finally monitoring strategies are evaluated. 

induced by  is a major cause of 
morbidity and mortality in  Despite a growing understanding of the
  and associated diseases in there 
 remains uncertainty regarding the best clinical management of these patients
These guidelines are offered to improve care by establishing consistent evidence based 
care Where evidence in the population is lacking literature from adult 
patients was reviewed and used where applicable and/or local consensus was used to form 
the recommendations These guidelines do not address 

A joint working group established by the  subgroup
 of the  Committee for Standards in  and the 
  Society has reviewed the available literature and made 
 recommendations for the diagnosis and management of 
 disorder in adult  This review details the 
 therapeutic options recommended including reduction in 
  Effective therapy should be 
 instituted before progressive disease results in declining performance status and 
  dysfunction. The goal of treatment should be a durable complete remission 
 with retention of  function with minimal

 represents a 
spectrum of  clinical diseases, from a benign 
like illness to a fulminant In the setting of 
 is an often fatal complication occurring 
relatively early after Risk factors for the development of are well 
established and include  depletion, and the use of 
as conditioning or treatment of versus  disease. 
Early recognition of  is particularly important in the  setting because  in 
these patients tends to be rapidly progressive. Familiarity with the clinical features of 
 and a heightened level of suspicion are critical for making the diagnosis. 
Surveillance techniques with 
may have a role in some high risk settings based therapies such as 
 either as treatment 
for  in high risk patients represent promising new directions in 
the treatment of this disease `

const NonKeyWords string = ` a acquire address adult after an and and/or applicable 
approaches are as associated atypical available based be because before benign best by 
care carry cause clinical committee complete complication conditioning consensus 
consistent critical current declining depending depletion despite details develop 
developing development diagnosis differential difficult directions discusses disease 
diseases disorder do dramatically due durable dysfunction early effective either 
especially established establishing evaluated eventually evidence factors familiarity 
fatal features finally for form from fulminant function gives goal group growing 
guidelines hampered has have heightened high however illness important improve improved 
in include includes including increased induced instituted introduction is joint lacking 
level like limited literature local made maintain major making management may minimal 
modalities monitoring morbidity mortality new not occurring of offered often on options 
or outcomes outlook part particularly patients performance population progressive 
promising rapidly recognition recommendations recommended reduction regarding relatively 
remains remission represent represents required results retention review reviewed risk 
role setting settings should society some spectrum standards status strategies subgroup 
substantially succumb such summarizes surveillance suspicion targeted techniques tends 
the therapeutic therapies therapy there these this those to towards treatment uncertainty 
understanding unspecific up use used versus vulnerability was well where which with 
working`

const keyWords string = ` cell antibod barr chemotherapy children donor ebv epstein 
graft haemato oncology hla immun infect interferon*a lympho
mono organ hodgkin pcr pediatric transplant prophylaxis ptld radiotherapy recipient 
rejection ris rituximab sct solid sot stem titer toxic virus
 `
