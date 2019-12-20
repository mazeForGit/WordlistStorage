package data

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"errors"
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
	
    requestUrl = GlobalConfig.WordListUrl + "/wordlist"
	// todo: faster, but include test
    //requestUrl = GlobalConfig.WordListUrl + "/words?testOnly=true&format=json"
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
	// todo: faster, but include test
    //json.Unmarshal(body, &GlobalWordList.Words)
	fmt.Println("got .. GlobalWordList.Words = " + strconv.Itoa(len(GlobalWordList.Words)))

	return nil
}
