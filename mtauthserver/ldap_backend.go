package mtauthserver

import (
	"crypto/tls"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	log "github.com/sirupsen/logrus"
	"os"
)

func ConnectLDAP() *ldap.Conn {
	cnf := LoadConfig()

	l, err := ldap.DialURL(cnf.LdapServer, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: true}))

	// Try server 2
	if err != nil {
		log.Errorf("Connect to server1 %s failed. Error:%v", cnf.LdapServer, err)
		l, err := ldap.DialURL(cnf.LdapServer2, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: true}))

		if err != nil {
			log.Errorf("Connect to server2 %s failed. Error:%v", cnf.LdapServer2, err)
		} else {
			log.Infof("Connect successful to server 2 %s.", cnf.LdapServer2)
		}
		return l
	} else {
		log.Infof("Connect successful to server 1 %s.", cnf.LdapServer)
	}
	return l
}

func checkAuthentication(ldapConn *ldap.Conn, username string, password string) bool {

	cnf := LoadConfig()

	err := ldapConn.Bind(cnf.BindDN, cnf.BindPassword)

	if err != nil {
		log.Printf("Error ldap bind: %v", err)
		return false
	}

	searchFilter := fmt.Sprintf(cnf.UserFilter, username)
	searchRequest := ldap.NewSearchRequest(cnf.SearchBase, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, searchFilter, []string{"dn"}, nil)

	sr, err := ldapConn.Search(searchRequest)

	if err != nil {
		log.Printf("Error search user: %v", err)
		return false
	}

	if len(sr.Entries) != 1 {
		log.Println("User not found")
		return false
	}

	userDN := sr.Entries[0].DN
	err = ldapConn.Bind(userDN, password)

	if err != nil {
		log.Printf("Ldap auth failed: %v", err)
		return false
	}

	return true
}

func findInheritedGroups(ldapConn *ldap.Conn, userDN string, allGroups map[string]bool) error {

	searchRequest := ldap.NewSearchRequest(
		userDN,
		ldap.ScopeBaseObject,
		ldap.NeverDerefAliases,
		0, 0, false,
		"(objectClass=*)",
		[]string{"memberOf"},
		nil,
	)

	sr, err := ldapConn.Search(searchRequest)
	if err != nil {
		return err
	}

	for _, entry := range sr.Entries {
		for _, groupDN := range entry.GetAttributeValues("memberOf") {
			allGroups[groupDN] = true
			if err := findInheritedGroups(ldapConn, groupDN, allGroups); err != nil {
				return err
			}
		}
	}

	return nil
}

func getGroupCNs(ldapConn *ldap.Conn, groupDNs []string) ([]string, error) {
	var groupCNs []string

	for _, groupDN := range groupDNs {
		// Suchanfrage, um die Gruppeninformationen zu erhalten
		searchRequest := ldap.NewSearchRequest(
			groupDN,
			ldap.ScopeBaseObject,
			ldap.NeverDerefAliases,
			0, 0, false,
			"(objectClass=*)",
			[]string{"cn"},
			nil,
		)

		sr, err := ldapConn.Search(searchRequest)
		if err != nil {
			return nil, err
		}

		if len(sr.Entries) > 0 {
			groupCN := sr.Entries[0].GetAttributeValue("cn")
			groupCNs = append(groupCNs, groupCN)
		}
	}

	return groupCNs, nil
}

func getGroupsofUser(ldapConn *ldap.Conn, username string) []string {

	cnf := LoadConfig()
	groups := []string{}
	groupsdn := []string{}

	err := ldapConn.Bind(cnf.BindDN, cnf.BindPassword)

	if err != nil {
		log.Printf("Error ldap bind: %v", err)
		os.Exit(1)
	}

	searchFilter := fmt.Sprintf(cnf.UserFilter, username)
	searchRequest := ldap.NewSearchRequest(cnf.SearchBase, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, searchFilter, []string{"dn"}, nil)

	sr, err := ldapConn.Search(searchRequest)

	if err != nil {
		log.Printf("Error search user: %v", err)
		os.Exit(1)
	}

	if len(sr.Entries) != 1 {
		log.Println("User not found")
		os.Exit(1)
	}

	userDN := sr.Entries[0].DN

	allGroups := make(map[string]bool)

	if err := findInheritedGroups(ldapConn, userDN, allGroups); err != nil {
		log.Errorf("Error searching groups: %v", err)
	}

	for groupDN := range allGroups {
		groupsdn = append(groupsdn, groupDN)
	}

	groups, err = getGroupCNs(ldapConn, groupsdn)

	if err != nil {
		log.Errorf("Error searching CNs from group DNs: %v", err)
	}

	return groups
}
