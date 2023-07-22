package schedulers

import (
	orgRepo "balkantask/database/org"
	userRepo "balkantask/database/user"
	"fmt"
	"time"
)

func deleteDeactivatedUser() {
	fmt.Println("Executing the Deactivated User Deletion at", time.Now())
	// Calculate the date 30 days ago from today
	threshold := time.Now().AddDate(0, 0, -30)

	// Find the users whose status is "DEACTIVATED" and were last updated 30 days ago
	users, err := userRepo.GetDeactivatedUserForThreshold(threshold)
	if err != nil {
		fmt.Println("Error getting deactivated users:", err)
		return
	}

	// Delete the selected users from the database
	err = userRepo.DeleteUsers(users)
	if err != nil {
		fmt.Println("Error deleting users:", err)
		return
	}
}

func markAccountDeleted() {
	fmt.Println("Marking DEACTIVATED Accounts as DELETED at", time.Now())
	// Calculate the date 5 days ago from today
	threshold := time.Now().AddDate(0, 0, -5)

	// Find the orgs whose status is "DEACTIVATED" and were last updated 5 days ago
	orgs, err := orgRepo.GetDeactivatedOrgsForThreshold(threshold)
	if err != nil {
		fmt.Println("Error getting deactivated users:", err)
		return
	}

	// Mark the orgs as "DELETED"
	for _, org := range orgs {
		org.AccountStatus = "DELETED"
	}

	_, err = orgRepo.UpdateOrgs(orgs)
	if err != nil {
		fmt.Println("Error updating orgs:", err)
		return
	}
}

func deleteAccountsData() {
	fmt.Println("Deleting DELETED Accounts data at", time.Now())
	// Calculate the date 45 days ago from today
	threshold := time.Now().AddDate(0, 0, -45)

	// Find the orgs whose status is "DELETED" and were last updated 45 days ago
	orgs, err := orgRepo.GetDeletedOrgsForThreshold(threshold)

	if err != nil {
		fmt.Println("Error getting deleted users:", err)
		return
	}

	// Delete the orgs from the database
	err = orgRepo.DeleteOrgs(orgs)
	if err != nil {
		fmt.Println("Error deleting orgs:", err)
		return
	}

}

func Scheduler() {
	for {
		now := time.Now()
		nextRun := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 30, 0, 0, time.Local)

		// Calculate the duration until the next run
		duration := nextRun.Sub(now)

		// Wait until the next run time
		time.Sleep(duration)

		go deleteDeactivatedUser()
		go markAccountDeleted()
		go deleteAccountsData()
	}
}
