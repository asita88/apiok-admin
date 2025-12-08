package services

import (
	"apiok-admin/app/packages"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-ldap/ldap/v3"
)

func LdapAuthenticate(username string, password string) (string, string, error) {
	conf := packages.GetConfig()
	if conf == nil {
		return "", "", fmt.Errorf("LDAP is not enabled")
	}

	confValue := reflect.ValueOf(conf).Elem()
	ldapField := confValue.FieldByName("Ldap")
	if !ldapField.IsValid() {
		return "", "", fmt.Errorf("LDAP configuration not found")
	}

	enabled := ldapField.FieldByName("Enabled").Bool()
	if !enabled {
		return "", "", fmt.Errorf("LDAP is not enabled")
	}

	host := ldapField.FieldByName("Host").String()
	baseDN := ldapField.FieldByName("BaseDN").String()
	if host == "" || baseDN == "" {
		return "", "", fmt.Errorf("LDAP configuration is incomplete")
	}

	port := int(ldapField.FieldByName("Port").Int())
	if port == 0 {
		port = 389
	}
	bindDN := ldapField.FieldByName("BindDN").String()
	bindPassword := ldapField.FieldByName("BindPassword").String()
	userFilter := ldapField.FieldByName("UserFilter").String()

	attributesField := ldapField.FieldByName("Attributes")
	nameAttr := attributesField.FieldByName("Name").String()
	emailAttr := attributesField.FieldByName("Email").String()

	if nameAttr == "" {
		nameAttr = "cn"
	}
	if emailAttr == "" {
		emailAttr = "mail"
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := ldap.Dial("tcp", addr)
	if err != nil {
		return "", "", fmt.Errorf("failed to connect to LDAP server: %v", err)
	}
	defer conn.Close()

	if bindDN != "" && bindPassword != "" {
		err = conn.Bind(bindDN, bindPassword)
		if err != nil {
			return "", "", fmt.Errorf("failed to bind with bind DN: %v", err)
		}
	}

	if userFilter == "" {
		userFilter = "(uid=%s)"
	}
	userFilter = strings.ReplaceAll(userFilter, "%s", ldap.EscapeFilter(username))

	searchRequest := ldap.NewSearchRequest(
		baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		userFilter,
		[]string{"dn", nameAttr, emailAttr},
		nil,
	)

	sr, err := conn.Search(searchRequest)
	if err != nil {
		return "", "", fmt.Errorf("failed to search user: %v", err)
	}

	if len(sr.Entries) == 0 {
		return "", "", fmt.Errorf("user not found")
	}

	if len(sr.Entries) > 1 {
		return "", "", fmt.Errorf("multiple users found")
	}

	userDN := sr.Entries[0].DN

	err = conn.Bind(userDN, password)
	if err != nil {
		return "", "", fmt.Errorf("invalid password")
	}

	name := ""
	email := ""

	if nameValue := sr.Entries[0].GetAttributeValue(nameAttr); nameValue != "" {
		name = nameValue
	} else {
		name = username
	}

	if emailValue := sr.Entries[0].GetAttributeValue(emailAttr); emailValue != "" {
		email = emailValue
	} else {
		email = username
	}

	return name, email, nil
}

func CheckUserAndPasswordWithLdap(email string, password string) error {
	conf := packages.GetConfig()
	if conf == nil {
		return CheckUserAndPassword(email, password)
	}

	confValue := reflect.ValueOf(conf).Elem()
	ldapField := confValue.FieldByName("Ldap")
	if ldapField.IsValid() && ldapField.FieldByName("Enabled").Bool() {
		_, _, err := LdapAuthenticate(email, password)
		if err != nil {
			return fmt.Errorf("LDAP authentication failed: %v", err)
		}
		return nil
	}

	return CheckUserAndPassword(email, password)
}

func GetUserEmailWithLdap(email string) string {
	conf := packages.GetConfig()
	if conf == nil {
		return email
	}

	confValue := reflect.ValueOf(conf).Elem()
	ldapField := confValue.FieldByName("Ldap")
	if ldapField.IsValid() && ldapField.FieldByName("Enabled").Bool() {
		_, ldapEmail, err := LdapAuthenticate(email, "")
		if err == nil && ldapEmail != "" {
			return ldapEmail
		}
	}
	return email
}
