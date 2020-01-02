package data

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"errors"
	"math/rand"
	"time"
)
type Config struct {
    WordListUrl  string	`json:"wordlisturl"`
}
type ResponseStatus struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}
type SessionStatus struct {
    SessionID int `json:"sid"`
    Count int `json:"count"`
    RequestExecution bool `json:"requestexecution"`
	PageToScan string `json:"pagetoscan"`
	DomainsAllowed string `json:"domainsallowed"`
	NumberLinksFound int `json:"numberlinksfound"`
	NumberLinksVisited int `json:"numberlinksvisited"`
	ExecutionStarted bool `json:"executionstarted"`
	ExecutionFinished bool `json:"executionfinished"`
	WordsScanned int `json:"wordsscanned"`
}
type Test struct {
    Name  string	`json:"name"`
    Category  string	`json:"category"`
	Words []int	`json:"words"`
}
type Word struct {
	Id  int		`json:"id"`
    Name  string	`json:"name"`
	Count int		`json:"count"`
	Occurance int	`json:"occurance"`
	Tests []Test	`json:"tests"`
}
type WordList struct {
	Session SessionStatus		`json:"session"`
	Words []Word				`json:"words"`
	Tests []Test				`json:"tests"`
}
type ByCategory struct {
	Name string					`json:"name"`
	ReferenceWordCount int		`json:"referencewordcount"`
	ReferenceCountSum int		`json:"referencecountsum"`
	ReferenceOccuranceSum int	`json:"referenceoccurancesum"`
	GlobalWordCount int			`json:"globalwordcount"`
	GlobalCountSum int			`json:"globalcountsum"`
	GlobalOccuranceSum int		`json:"globaloccurancesum"`
	LocalWordCount int			`json:"localwordcount"`
	LocalOccuranceSum int		`json:"localoccurancesum"`
}
type Result struct {
	Test string				`json:"test"`
	Category []ByCategory
}

var GlobalConfig Config
var GlobalWordList WordList
var GlobalWordListStorage []WordList
var GlobalWordListResult Result


func GetDomains() []SessionStatus{
	fmt.Println("GetDomains")
	
	var ss []SessionStatus
	for i := 0; i < len(GlobalWordListStorage); i++ {
		ss = append(ss, GlobalWordListStorage[i].Session)
	}
	
	return ss
}
func GetWordList(test string, size string) ([]Word, error) {
	fmt.Println("GetWordList .. test = " + test + ", size = " + size)
	
	var wl []Word
	var maxWordsInTest int = 10
		
	if size == "short" {
		maxWordsInTest = 20
	} else if size == "medium" {
		maxWordsInTest = 40
	} else if size == "long" {
		maxWordsInTest = 80
	}
	fmt.Println("GetWordList .. maxWordsInTest = " + strconv.Itoa(maxWordsInTest))
	
	// find out how many categories are in test
	var maxCategories = 0
	for k := 0; k < len(GlobalWordList.Tests); k++ {
		//fmt.Println("GetWordList .. GlobalWordList.Tests[k].Name = " + GlobalWordList.Tests[k].Name)
		if GlobalWordList.Tests[k].Name == test {
			maxCategories++
		}
	}
	fmt.Println("GetWordList .. maxCategories = " + strconv.Itoa(maxCategories))
	
	if maxCategories == 0 {
		return nil, errors.New("no test category")
	}
	
	var wordsPerTestCategory int = maxWordsInTest / maxCategories
	fmt.Println("GetWordList .. wordsPerTestCategory = " + strconv.Itoa(wordsPerTestCategory))
	
	s1 := rand.NewSource(time.Now().UnixNano())
    r1 := rand.New(s1)
	
	// random over all without test category
	//for i := 0; i < max; i++ {
	//	j := r1.Intn(len(GlobalWordList.Words) - 1)
	//	w := GlobalWordList.Words[j]
	//	w.Tests = nil
	//	w.Count = 0
	//	w.Occurance = 0
	//	wl = append(wl, w)
	//} 
	
	// random over test category
	for k := 0; k < len(GlobalWordList.Tests); k++ {
		if GlobalWordList.Tests[k].Name == test {
			fmt.Println("GetWordList .. Test.Name = " + GlobalWordList.Tests[k].Name + ", Test.Category = " + GlobalWordList.Tests[k].Category + ", Test.Words.length = " + strconv.Itoa(len(GlobalWordList.Tests[k].Words)))
			
			for i := 0; i < wordsPerTestCategory; i++ {
				j := r1.Intn(len(GlobalWordList.Tests[k].Words) - 1)
				rw := GlobalWordList.Tests[k].Words[j]
				w := GlobalWordList.Words[rw]
				w.Tests = nil
				w.Count = 0
				w.Occurance = 0
				wl = append(wl, w)
			} 
		}
	}
	
	return wl, nil
}
func RebuildWordListResult() {
	fmt.Println("RebuildWordListResult")
	
	//fmt.Println("len(GlobalWordListStorage) = " + strconv.Itoa(len(GlobalWordListStorage)))
	for i := 0; i < len(GlobalWordListStorage); i++ {
		//fmt.Println(" len(GlobalWordListStorage[i].Words) = " + strconv.Itoa(len(GlobalWordListStorage[i].Words)))
		for j := 0; j < len(GlobalWordListStorage[i].Words); j++ {
			for k := 0; k < len(GlobalWordList.Words); k++ {
				if GlobalWordListStorage[i].Words[j].Name == GlobalWordList.Words[k].Name {
					//fmt.Println("name = " + GlobalWordListStorage[i].Words[j].Name + ", occurance = " + strconv.Itoa(GlobalWordListStorage[i].Words[j].Occurance))
					GlobalWordList.Words[k].Count++
					GlobalWordList.Words[k].Occurance += GlobalWordListStorage[i].Words[j].Occurance
					break
				}
			}
		}
	}	
}
func ResultOnSession(test string, domain string) ([]Word, error) {
	fmt.Println("ResultOnSession .. domain = " + domain)
	
	sData, err := GetWordListFromStorage(domain)
	if err != nil {
		return nil, errors.New("no item")
	}
	gwl := GetWordsList(test)
	lwl := sData.Words
	
	for i := 0; i < len(lwl); i++ {
		for j := 0; j < len(gwl); j++ {
			if lwl[i].Name == gwl[j].Name {
				if gwl[j].Occurance > 0 && gwl[j].Count > 0 {
					lwl[i].Id = gwl[j].Id
					lwl[i].Count = gwl[j].Occurance / gwl[j].Count
				}
			}
		}
	}
						
	return lwl, nil
}
func ResultOnSessionByTest(test string, domain string) (Result, error) {
	fmt.Println("ResultOnSession .. test = " + test + ", domain = " + domain)
	
	gwl := GetWordsList(test)
	PrepareResultsBasedOnTest(test, gwl)
	sData, err := GetWordListFromStorage(domain)
	if err != nil {
		return GlobalWordListResult, errors.New("no item")
	}
	
	//fmt.Println(sData)
	lwl := sData.Words
	
	//fmt.Println("len(GlobalWordListResult.Category) = " + strconv.Itoa(len(GlobalWordListResult.Category)))
	//fmt.Println("len(gwl) =" + strconv.Itoa(len(gwl)))
	//fmt.Println("len(lwl) =" + strconv.Itoa(len(lwl)))
	
	for i := 0; i < len(gwl); i++ {
		for ii := 0; ii < len(gwl[i].Tests); ii++ {
			if gwl[i].Tests[ii].Name == test {
				for j := 0; j < len(lwl); j++ {
					if gwl[i].Name == lwl[j].Name {
						for k := 0; k < len(GlobalWordListResult.Category); k++ {
							if GlobalWordListResult.Category[k].Name == gwl[i].Tests[ii].Category {
								GlobalWordListResult.Category[k].LocalWordCount++
								GlobalWordListResult.Category[k].LocalOccuranceSum += lwl[j].Occurance
								
								break
							}
						}
					}
				}
				for k := 0; k < len(GlobalWordListResult.Category); k++ {
					if GlobalWordListResult.Category[k].Name == gwl[i].Tests[ii].Category {			
						if gwl[i].Occurance > 0 {
							GlobalWordListResult.Category[k].GlobalCountSum += gwl[i].Count
							GlobalWordListResult.Category[k].GlobalWordCount++
							GlobalWordListResult.Category[k].GlobalOccuranceSum += gwl[i].Occurance
						}
						break
					}
				}
			}
		} 
	}
	
	return GlobalWordListResult, nil
}
func PrepareResultsBasedOnTest(test string, gwl []Word) {
	//fmt.Println("PrepareResultsBasedOnTest .. test = " + test)

	GlobalWordListResult.Test = test
	GlobalWordListResult.Category = GlobalWordListResult.Category[0:0]
	for i := 0; i < len(GlobalWordList.Tests); i++ {
		//fmt.Println("GlobalWordList.Tests[i].Name = " + GlobalWordList.Tests[i].Name)
		if GlobalWordList.Tests[i].Name == test {
			var c = ByCategory {
				Name: GlobalWordList.Tests[i].Category,
				ReferenceWordCount: 0,
				ReferenceCountSum: 0,
				ReferenceOccuranceSum: 0,
				GlobalCountSum: 0,
				GlobalWordCount: 0,
				GlobalOccuranceSum: 0,
				LocalWordCount: 0,
				LocalOccuranceSum: 0,
			}
			GlobalWordListResult.Category = append(GlobalWordListResult.Category, c)
		}
	}
	// init with reference from GlobalWordList
	for i := 0; i < len(gwl); i++ {
		for ii := 0; ii < len(gwl[i].Tests); ii++ {
			if gwl[i].Tests[ii].Name == test {
				
						for k := 0; k < len(GlobalWordListResult.Category); k++ {
							if GlobalWordListResult.Category[k].Name == gwl[i].Tests[ii].Category {
								GlobalWordListResult.Category[k].ReferenceWordCount++ 
								GlobalWordListResult.Category[k].ReferenceCountSum += gwl[i].Count 
								GlobalWordListResult.Category[k].ReferenceOccuranceSum += gwl[i].Occurance 
								break
							}
						}
			}
		} 
	}
	
	//fmt.Println(GlobalWordListResult)
}
func GetWordsList(test string) []Word {
	fmt.Println("GetWordsList .. test = " + test)
	var wl []Word
	
	fmt.Println("GetWordsList .. len(GlobalWordList.Words) = " + strconv.Itoa(len(GlobalWordList.Words)))
	
	for i := 0; i < len(GlobalWordList.Words); i++ {
		if GlobalWordList.Words[i].Tests != nil {
			for j := 0; j < len(GlobalWordList.Words[i].Tests); j++ {
				if test == "" || GlobalWordList.Words[i].Tests[j].Name == test {
					wl = append(wl, GlobalWordList.Words[i])
					//fmt.Println(GlobalWordList.Words[i])
					break;
				}
			}
		} 
	}
	
	fmt.Println("GetWordsList .. len(wl) = " + strconv.Itoa(len(wl)))
	return wl
}
func GetWordListFromStorage(domain string) (WordList, error) {
	fmt.Println("GetWordListFromStorage .. domain = " + domain)
    var wl WordList
	
	for i := 0; i < len(GlobalWordListStorage); i++ {
		if GlobalWordListStorage[i].Session.DomainsAllowed == domain {
			return GlobalWordListStorage[i], nil
		}
	}

    return wl, errors.New("no item")
}
func AddWordsToStorage(test string, domain string, wrds []Word) {
	fmt.Println("AddWordsToStorage.. test = " + test + ", domain = " + domain)
	
	if len(wrds) == 0 {
		return
	}
	// insert into storage
	// replace if present
	
	//fmt.Println(wrds)
	var wl WordList
	wl.Words = wrds
	wl.Session = SessionStatus { SessionID: 0, Count: 0, RequestExecution: false, PageToScan: test, DomainsAllowed: domain, NumberLinksFound: 0, NumberLinksVisited: 0, ExecutionStarted: false, ExecutionFinished: true, WordsScanned: 0 }
	
	AddWordListToStorage(wl)
}
func AddWordListToStorage(wl WordList) {
	fmt.Println("AddWordListToStorage.. len(wl.Words) = " + strconv.Itoa(len(wl.Words)))
	
	if len(wl.Words) == 0 {
		return
	}
	// insert into storage
	// replace if present
	for i := 0; i < len(GlobalWordListStorage); i++ {
		if GlobalWordListStorage[i].Session.DomainsAllowed == wl.Session.DomainsAllowed {
			// remove results
			RemoveWordListResultsFromGlobalWordList(GlobalWordListStorage[i].Words)
			
			// Remove the element at index i from wl
			copy(GlobalWordListStorage[i:], GlobalWordListStorage[i+1:]) 					// Shift a[i+1:] left one index.
			GlobalWordListStorage = GlobalWordListStorage[:len(GlobalWordListStorage)-1]  	// Truncate slice.
			break;
		}
	}
	GlobalWordListStorage = append(GlobalWordListStorage, wl)
	
	// add at GlobalWordList
	AddWordListResultsToGlobalWordList(wl.Words)
}
func RemoveWordListResultsFromGlobalWordList(wl []Word) {
	for i := 0; i < len(wl); i++ {
				for j := 0; j < len(GlobalWordList.Words); j++ {
					if wl[i].Name == GlobalWordList.Words[j].Name {
								GlobalWordList.Words[j].Occurance -= wl[i].Occurance
								GlobalWordList.Words[j].Count--
								break
					}
				}
	}
}
func AddWordListResultsToGlobalWordList(wl []Word) {
	fmt.Println("AddWordListResultsToGlobalWordList .. len(wl) = " + strconv.Itoa(len(wl)))
	
	for i := 0; i < len(wl); i++ {
				for j := 0; j < len(GlobalWordList.Words); j++ {
					if wl[i].Name == GlobalWordList.Words[j].Name {
								GlobalWordList.Words[j].Occurance += wl[i].Occurance
								GlobalWordList.Words[j].Count++
								break
					}
				}
	}
}

func ReadGlobalWordlistFromRemote() error {
	fmt.Println("ReadGlobalWordlist")
	fmt.Println("have .. GlobalWordlist.Words = " + strconv.Itoa(len(GlobalWordList.Words)))
	
    var err error
	var resp *http.Response
	var body []byte
	var requestUrl string = ""
	
    requestUrl = GlobalConfig.WordListUrl + "/wordlist?testOnly=true"
	fmt.Println("connect to wordlist and get words with tests = " + requestUrl)
    resp, err = http.Get(requestUrl)
    if err != nil {
        return err
    }

    defer resp.Body.Close()
    body, err = ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }
    
    json.Unmarshal(body, &GlobalWordList)
	fmt.Println("got .. GlobalWordList.Words = " + strconv.Itoa(len(GlobalWordList.Words)))

	PrepareGlobalWordlistForTest()
	
	return nil
}
func PrepareGlobalWordlistForTest() {
	fmt.Println(PrepareGlobalWordlistForTest)
	
	// 
	for i := 0; i < len(GlobalWordList.Words); i++ {
		for j := 0; j < len(GlobalWordList.Words[i].Tests); j++ {
			for k := 0; k < len(GlobalWordList.Tests); k++ {
				if GlobalWordList.Words[i].Tests[j].Name == GlobalWordList.Tests[k].Name && GlobalWordList.Words[i].Tests[j].Category == GlobalWordList.Tests[k].Category {
					GlobalWordList.Tests[k].Words = append(GlobalWordList.Tests[k].Words, i)
				}
			}
		} 
	}
	
	//fmt.Println(GlobalWordList.Tests)
}