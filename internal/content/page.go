package content

//
// Page
//
//
//
// type Page struct {
// 	Title       string
// 	Name        string
// 	contentPath string
// }

//
//
// PageIndex
//
//
// type PageIndex map[string]*Page
//
// func (p PageIndex) Routes() ([]string, error) {
// 	return []string{}, nil
// }
//
// func LoadPages(contentDir string) (PageIndex, error) {
// 	root := filepath.Join(contentDir, "pages")
// 	paths := make([]string, 0)
// 	filepath.Walk(root,
// 		func(path string, info os.FileInfo, err error) error {
// 			if err != nil || info.IsDir() {
// 				return nil
// 			}
// 			if strings.HasSuffix(path, ".md") {
// 				paths = append(paths, path)
// 			}
// 			return nil
// 		})
//
// 	log.Printf("LoadPages: paths=%v", paths)
//
// 	return nil, nil
// }
