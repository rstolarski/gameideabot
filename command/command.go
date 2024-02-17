package command

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func Add(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.Contains(m.Content, "!add") && m.Content[0:1] == "!" {
		content := strings.Join(strings.Split(m.Content, " ")[1:], " ")

		file, err := os.OpenFile("list.txt", os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf("Failed opening file: %s", err)
		}
		defer file.Close()

		_, err = file.WriteString("\n" + content)
		if err != nil {
			log.Fatalf("Failed writing to file: %s", err)
		}

		log.Printf("Added new entry: %s", content)
		s.ChannelMessageSendReply(m.ChannelID, "Ok, dodałem nowe hasło - "+content, &discordgo.MessageReference{MessageID: m.ID})
	}
}

func Del(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.Contains(m.Content, "!del") && m.Content[0:1] == "!" {
		entry_to_del := strings.Join(strings.Split(m.Content, " ")[1:], " ")
		log.Println("Deleting entry:", entry_to_del)

		file, err := os.Open("list.txt")
		if err != nil {
			log.Fatalf("Failed opening file: %s", err)
		}

		defer file.Close()

		var lines []string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		if len(lines) == 0 {
			log.Fatalln("File is empty")
		}

		found := false
		for i, v := range lines {
			if v == entry_to_del {
				lines = append(lines[:i], lines[i+1:]...)
				found = true
				break
			}
		}

		if !found {
			log.Printf("Failed to find entry: %s", entry_to_del)

			s.ChannelMessageSendReply(m.ChannelID, "Nie mam tego na liście. Literówka?", &discordgo.MessageReference{MessageID: m.ID})
		} else {
			err = os.WriteFile("list.txt", []byte{}, 0664)
			if err != nil {
				log.Fatalf("Failed writing to file: %s", err)
			}

			file, err = os.OpenFile("list.txt", os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				log.Fatalf("Failed opening file: %s", err)
			}
			defer file.Close()

			file_content := strings.Join(lines, "\n")

			_, err = file.WriteString(file_content)
			if err != nil {
				log.Fatalf("Failed writing to file: %s", err)
			}

			log.Printf("%s has been deleted.", entry_to_del)
			s.ChannelMessageSendReply(m.ChannelID, entry_to_del+" wyrzucone z listy", &discordgo.MessageReference{MessageID: m.ID})
		}
	}

}

func List(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.Contains(m.Content, "!ls") && m.Content[0:1] == "!" {
		file, err := os.Open("list.txt")
		if err != nil {
			log.Fatalln(err)
		}

		defer file.Close()

		var lines []string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		log.Printf("Current list:\n%s", strings.Join(lines, "\n"))
		s.ChannelMessageSendReply(m.ChannelID, "Obecnie na liście mamy:\n"+strings.Join(lines, "\n"), &discordgo.MessageReference{MessageID: m.ID})
	}
}

func Random(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.Contains(m.Content, "!random") {
		number := 1
		var err error

		if r, _ := strings.CutPrefix(strings.TrimSpace(m.Content), "!random"); r != "" {
			number, err = strconv.Atoi(strings.Split(m.Content, " ")[1])
			if err != nil {
				panic(err)
			}
		}

		file, err := os.Open("list.txt")
		if err != nil {
			log.Fatalln(err)
		}
		defer file.Close()

		var lines []string
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(lines), func(i, j int) { lines[i], lines[j] = lines[j], lines[i] })

		msg := strings.Join(lines[0:number], ", ")
		s.ChannelMessageSendReply(m.ChannelID, msg, &discordgo.MessageReference{MessageID: m.ID})
	}
}
