package processing

import (
  "fmt"
  "testing"
)

func BenchmarkTokenizer(b *testing.B) {
  for n := 0 ; n < b.N; n ++ {
    SimpleTokenizer(`
oundCloud is an online audio distribution platform based in Berlin, Germany that enables its users to upload, record, promote and share their originally-created sounds. In July 2013, it had 40 million registered users and 200 million listeners[3]

Contents  [hide] 
1 History
2 Features
3 Paid subscription features
4 Recognition
5 Criticism
6 Blocking
7 References
8 External links

History[edit]
SoundCloud was originally started in Stockholm, Sweden, but it was established in Berlin in August 2007 by Swedish sound designer Alex Ljung and Swedish artist Eric Wahlforss. It had the intention of allowing musicians to share recordings with each other, but it later transformed into a full publishing tool which also allowed musicians to distribute their music tracks.[4]

A few months after it began operating, SoundCloud began to challenge the dominance of Myspace as a platform for musicians to distribute their music by allowing recording artists to interact more nimbly with their fans.[4]

In a 2009 interview with Wired, co-founder Alex Ljung said:

We both came from backgrounds connected to music, and it was just really, really annoying for us to collaborate with people on music—I mean simple collaboration, just sending tracks to other people in a private setting, getting some feedback from them, and having a conversation about that piece of music. In the same way that we’d be using Flickr for our photos, and Vimeo for our videos, we didn't have that kind of platform for our music.[4]

In April 2009, SoundCloud received €2.5 million Series A funding from Doughty Hanson Technology Ventures.[5] By May 2010, SoundCloud announced it had one million subscribers.[5] In January 2011, it was confirmed that SoundCloud had raised a $10 million Series B funding round from Union Square Ventures and Index Ventures. On 15 June 2011, SoundCloud announced they had five million registered users, and investments from Ashton Kutcher and Guy Oseary's A-Grade Fund. On 23 January 2012, SoundCloud announced on their blog that they had 10 million registered users. A story wheel was created for the occasion, which can be found on the SoundCloud blog.[6] In December 2012, a new SoundCloud layout was released to the general public. Among the many new features was the ability to continue playback of a track whilst navigating around the site, and the ability to read comments without them obscuring the waveform.

In March 2014, it was reported SoundCloud was in talks with major music labels regarding licensing due to copyrighted material appearing on the platform. This is in attempt to avoid the situation at Google and Youtube who are forced to handle a large number of takedown notices.[7]

Features[edit]
One of the key features of SoundCloud is that it lets artists upload their music with a distinctive URL. By allowing sound files to be embedded anywhere, SoundCloud can be combined with Twitter and Facebook to let members reach their audience better. This contrasts with MySpace, which hosts music only on the MySpace site.[4]

Registered SoundCloud users have the power to listen to as much content as they wish, as well as the ability to download up to 100 songs from the site. SoundCloud also allows its musician users to upload up to 120 minutes of audio to their profile. All of these features are free of charge and are available to all SoundCloud users as soon as they have a registered profile with the site.[8]

Astonishingly, an average of 12 hours of audio is uploaded to the site every minute, and half of SoundCloud’s content is split between small-time musicians and mainstream superstars.[9]

SoundCloud distributes music using widgets and apps.[5] Users can place the widget on their own websites or blogs and then SoundCloud will automatically tweet every track uploaded.[4] SoundCloud has an API that allows other applications or smartphones to upload or download music and sound files.[4]

This API has been integrated into several applications, most notably GarageBand, Logic Pro, and PreSonus Studio One DAW.[10] The API is also integrated in music finders, including SoundYouNeed.[11] Users may also download a music with creative commons licence by this API.

SoundCloud depicts audio tracks graphically as waveforms and allows users to comment on specific parts of the track (also known as timed comments). These comments are displayed while listening to the part of the sound they are referring to. Other standard features include reposts, playlists (previously known as "sets"),[citation needed] followers, and complimentary digital downloads.[12]

SoundCloud also provides users with the ability to create and join groups that provide a common space for content to be collected and shared.

Of the songs that are uploaded, over half of them are played within the first 30 minutes, and 90% of all uploaded tracks receive a listen from at least one user.[13]

Paid subscription features[edit]
SoundCloud offers additional features to users with paid subscriptions. If a user wishes to upload content that exceeds the initial free 2 hours, they may subscribe to a $38/year plan for up to 4 hours of content, or $130/year plan for unlimited uploads.[14] Such users are given more hosting space and may distribute their tracks or recordings to more groups and users, create sets of recordings, and more thoroughly track the statistics for each of their tracks. Additional statistic data are unlocked depending on which subscription the user has chosen, including the number of listens per track per user and the originating country of individual listens.[15]

SoundCloud also releases applications for popular mobile operating systems like iOS[16] and Android[17]

Recognition[edit]
SoundCloud won the Schroders Innovation Award at the 2011 European Tech Tour Awards Dinner.[18][19]

Criticism[edit]
As SoundCloud has grown and expanded beyond its initial user base primarily consisting of grassroots musicians, some original users have complained that it is losing its fidelity to artists in an attempt to appeal to the masses, perhaps in preparation for public sale.[20] Such criticism particularly followed the launching of a revamped website in 2013 that was heavily reconfigured to be more amenable to listeners—at the expense of artists, some claimed. CEO Ljung responded[citation needed] that while he would take these criticisms into consideration, these listener-friendly changes would likely attract many new users.

Blocking[edit]
Access to SoundCloud website has been blocked by the Government of Turkey on 24 January 2014.[21][22][23][24] A user named "haramzadeler" ("illegal ones" in Turkish) uploaded a total of 7 secretly recorded phone call tapes which reveal private conversations between Recep Tayyip Erdoğan (the Turkish Prime Minister) and others, including: Erdoğan Bayraktar, local politicians, some businessmen, PM's daughter Sümeyye Erdoğan and his son Bilal Erdoğan. Linked to the 2013 corruption scandal in Turkey, sound recordings revealed some conversations of illegal activity and possible bribery, mainly about building permit of luxurious villas on protected cultural heritage sites in Urla, İzmir.[25]

Opposition party Cumhuriyet Halk Partisi submitted a parliamentary question to TBMM concerning the issue and the questionnaire particularly asks for reasons of banning SoundCloud services without any proper case and/or reason.[26][27]
    `, 128)
  }
}

func TestSimpleTokenizer(t *testing.T) {
  ref := []string{"See", "John", "Markoff", "Apple", "Adopts", "Open", "Source", "for", "its", "Server", "Computers", "New", "York", "Times"}
  tokenized := SimpleTokenizer(`
      #See John Markoff, “Apple Adopts ‘Open Source’ for its Server Computers, ‘New York Times’”
     `, 128)

  if fmt.Sprintf("%v",tokenized) != fmt.Sprintf("%v",ref) {
    t.Errorf("expected %v == %v", tokenized, ref)
  }
}

func TestSimpleTokenizerApostrophe(t *testing.T) {
  ref := []string{"There's","not","a","problem","that","I","can't","fix"}
  tokenized := SimpleTokenizer(`
     'There's not a problem that I can't fix'
  `, 128)
  if fmt.Sprintf("%v",tokenized) != fmt.Sprintf("%v",ref) {
    t.Errorf("expected %v == %v", tokenized, ref)
  }
}

// ignoring japanese, chinese, korean for now, sorry
func TestSimpleTokenizerCJK(t *testing.T) {
  ref := []string{"宮","崎","駿"}
  tokenized := SimpleTokenizer(`
     宮崎 駿
  `, 128)
  if fmt.Sprintf("%v",tokenized) != fmt.Sprintf("%v",ref) {
    t.Errorf("expected %v == %v", tokenized, ref)
  }
}

func TestSimpleTokenizerEmptyString(t *testing.T) {
  tokenized := SimpleTokenizer("", 128)
  if len(tokenized) != 0 {
    t.Errorf("expected empty array, but got %v", tokenized)
  }
}

func TestSimpleTokenizerCap(t *testing.T) {
  ref := []string{"not","a","that","I","can't","fix"}
  tokenized := SimpleTokenizer(`
     'There's not a problem that I can't fix'
  `, 4)
  if fmt.Sprintf("%v",tokenized) != fmt.Sprintf("%v",ref) {
    t.Errorf("expected %v == %v", tokenized, ref)
  }
}

func TestLowercaseFilter(t *testing.T) {
  src := []string{"Compare", "FOO"}
  ref := []string{"compare", "foo"}
  result := LowercaseFilter(src)
  if fmt.Sprintf("%v",result) != fmt.Sprintf("%v",ref) {
    t.Errorf("expected %v == %v", result, ref)
  }
}

func TestStopWordsFilter(t *testing.T) {
  filter := CreateStopWordsFilter([]string{"foo","bar","baz"})
  ref := []string{"quick","brown","fox","jumps","over","lazy","dog"}
  result := filter([]string{"foo","quick","bar","baz","brown","foo","fox","jumps","over","bar","lazy","dog","baz"})
  if fmt.Sprintf("%v",result) != fmt.Sprintf("%v",ref) {
    t.Errorf("expected %v == %v", result, ref)
  }
}
