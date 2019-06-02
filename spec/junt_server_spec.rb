require 'airborne'

describe 'Junt' do
  before(:all) do
    ENV['JUNT_DB'] = "./junt_test.db"
    # if ENV['JUNT_DB'].nil? 
    #   $stderr.puts("JUNT_DB env var must be specified.")
    #   exit 1
    # end
    # guarantee we're working with a fresh db
    `rm #{ENV['JUNT_DB']}`
    # start the junt server 
    @@junt_server_pid = fork do
      exec "./junt"
    end
    puts "about to sleep..."
    # Process.detach(junt_server)
    # give the server a couple seconds to create the db
    # and start listening
    sleep(2) 


    @@base_url = "http://localhost:8123/v1"
    @@default_headers = {} # uses application/json content type by default
  end
  after(:all) do
    # Kill the server
    Process.kill("HUP", @@junt_server_pid)
    # remove the test db
    # `rm #{ENV['JUNT_DB']}`
  end
  describe "initial lists" do
    it "should initially have no companies" do
      get "#{@@base_url}/companies" # should return empty array
      expect(json_body.size).to(eq(0))
    end
    it "should initially have no events" do
      get "#{@@base_url}/events" # should return empty array
      expect(json_body.size).to(eq(0))
    end
    it "should initially have no followups" do
      get "#{@@base_url}/followups" # should return empty array
      expect(json_body.size).to(eq(0))
    end
    it "should initially have no homeworks" do
      get "#{@@base_url}/homeworks" # should return empty array
      expect(json_body.size).to(eq(0))
    end
    it "should initially have no interviews" do
      get "#{@@base_url}/interviews" # should return empty array
      expect(json_body.size).to(eq(0))
    end
    it "should initially have no jobs" do
      get "#{@@base_url}/jobs" # should return empty array
      expect(json_body.size).to(eq(0))
    end
    it "should initially have no offers" do
      get "#{@@base_url}/offers" # should return empty array
      expect(json_body.size).to(eq(0))
    end
    it "should initially have no people" do
      get "#{@@base_url}/people" # should return empty array
      expect(json_body.size).to(eq(0))
    end
    it "should initially have no status_changes" do
      get "#{@@base_url}/status_changes" # should return empty array
      expect(json_body.size).to(eq(0))
    end
    it "should initially have no thanks_emails" do
      get "#{@@base_url}/thanks_emails" # should return empty array
      expect(json_body.size).to(eq(0))
    end
  end
  describe "creation" do
    it "should be able to create a company" do
      new_company_json = {
        name: "test co",
        location: "somewhere nice",
        url: "https://example.com/test_co",
        note: "test co has 100% coverage"
      }
      post "#{@@base_url}/companies", new_company_json, @@default_headers
      expect_status(200)
      expect_json_types(status: :string, id: :int)
      expect(json_body[:status]).to(eq("SUCCESS"))
      @@company_1_id = json_body[:id]
    end
    it "should be able to create a job" do
      new_job_json = {
        company_id: @@company_1_id,
        job_title: "chief tester",
        posting_url: "https://example.com/test_co/chief_tester_job",
        source: "we work remotely",
        referred_by: "nobody",
        salary_range: "unknown",
        application_method: "online",
        note: "a _markdown_ note",
        start_date: nil
      }
      post "#{@@base_url}/jobs", new_job_json, @@default_headers
      expect_status(200)
      expect_json_types(status: :string, id: :int)
      expect(json_body[:status]).to(eq("SUCCESS"))
      @@job_1_id = json_body[:id]
    end
    it "should be able to create a new person" do
      new_person_json = {
        name: "mary smith",
        email: 'msmith@@example.com',
        phone: '+1 555-555-5555',
        note: 'another _markdown_ note',
      }
      post "#{@@base_url}/people", new_person_json, @@default_headers
      expect_status(200)
      expect_json_types(status: :string, id: :int)
      expect(json_body[:status]).to(eq("SUCCESS"))
      @@person_1_id = json_body[:id]
    end
    describe "of events" do
      it "should be able to create a new followup" do
        new_followup_json = {
          job_id: @@job_1_id,
          person_ids: [@@person_1_id],
          note: "a _markdown_ note"
        }
        post "#{@@base_url}/followups", new_followup_json, @@default_headers
        expect_status(200)
        expect_json_types(status: :string, id: :int)
        expect(json_body[:status]).to(eq("SUCCESS"))
        @@followup_1_id = json_body[:id]
      end
      it "should be able to create a new homework" do
        new_homework_json = {
          job_id: @@job_1_id,
          due_date: "2019-06-02T18:47:44+00:00",
          note: "a _markdown_ note"
        }
        post "#{@@base_url}/homeworks", new_homework_json, @@default_headers
        expect_status(200)
        expect_json_types(status: :string, id: :int)
        expect(json_body[:status]).to(eq("SUCCESS"))
        @@homework_1_id = json_body[:id]
      end
      it "should be able to create a new interview" do
        new_interview_json = {
          job_id: @@job_1_id,
          person_ids: [@@person_1_id],
          note: "a _markdown_ note",
          scheduled_at: "2019-06-02T18:47:44+00:00",
          rating: "😃",
          type: "technical"
        }
        post "#{@@base_url}/interviews", new_interview_json, @@default_headers
        expect_status(200)
        expect_json_types(status: :string, id: :int)
        expect(json_body[:status]).to(eq("SUCCESS"))
        @@interview_1_id = json_body[:id]
      end
      it "should be able to create a new offer" do
        new_offer_json = {
          job_id: @@job_1_id,
          status: "received",
          note: "a _markdown_ note"
        }
        post "#{@@base_url}/offers", new_offer_json, @@default_headers
        expect_status(200)
        expect_json_types(status: :string, id: :int)
        expect(json_body[:status]).to(eq("SUCCESS"))
        @@offer_1_id = json_body[:id]
      end
      it "should be able to create a new status_change" do
        new_status_change_json = {
          job_id: @@job_1_id,
          from: nil, 
          # initial status change for a job should be nil->something
          to: "applied",
          note: "a _markdown_ note"
        }
        post "#{@@base_url}/status_changes", new_status_change_json, @@default_headers
        expect_status(200)
        expect_json_types(status: :string, id: :int)
        expect(json_body[:status]).to(eq("SUCCESS"))
        @@status_change_1_id = json_body[:id]
      end
      it "should be able to create a new thanks_email" do
        new_thanks_email_json = {
          job_id: @@job_1_id,
          person_ids: [@@person_1_id],
          note: "a _markdown_ note" 
        }
        post "#{@@base_url}/thanks_emails", new_thanks_email_json, @@default_headers
        expect_status(200)
        expect_json_types(status: :string, id: :int)
        expect(json_body[:status]).to(eq("SUCCESS"))
        @@thanks_email_1_id = json_body[:id]
      end
    end # end events creation
    
  end # end creation
  describe "showing objects" do
    it "should be able to show a company" do
      get "#{@@base_url}/companies/#{@@company_1_id}"
      expect_status(200)
      expect_json_types(
        id: :int,
        job_ids: :array_of_ints,
        person_ids: :array_of_ints,
        name: :string,
        location: :string_or_null,
        url: :string_or_null,
        note: :string_or_null,
        created_at: :date,
        updated_at: :date
      )
    end
    it "should be able to show a job" do
      get "#{@@base_url}/jobs/#{@@job_1_id}"
      expect_status(200)
      expect_json_types(
        id: :int,
        company_id: :int,
        person_ids: :array_of_ints,
        job_title: :string,
        posting_url: :string_or_null,
        source: :string_or_null,
        referred_by: :string_or_null,
        salary_range: :string_or_null,
        application_method: :string_or_null,
        note: :string_or_null,
        start_date: :date_or_null,
        created_at: :date,
        updated_at: :date
      )
    end
    it "should be able to show a person" do
      get "#{@@base_url}/people/#{@@person_1_id}"
      expect_status(200)
      expect_json_types(
        id: :int,
        job_ids: :array_of_ints,
        company_id: :int,
        name: :string,
        email: :string_or_null,
        phone: :string_or_null,
        created_at: :date,
        updated_at: :date
      )
    end
    # events....
    it "should be able to show a followup" do
      get "#{@@base_url}/followups/#{@@followup_1_id}"
      expect_status(200)
      expect_json_types(
        id: :int,
        job_id: :int,
        note: :string_or_null,
        person_ids: :array_of_ints,
        created_at: :date,
        updated_at: :date
      )
    end
    it "should be able to show a homework" do
      get "#{@@base_url}/homeworks/#{@@homework_1_id}"
      expect_status(200)
      expect_json_types(
        id: :int,
        job_id: :int,
        note: :string_or_null,
        due_date: :date,
        created_at: :date,
        updated_at: :date
      )
    end
    it "should be able to show a interview" do
      get "#{@@base_url}/interviews/#{@@interview_1_id}"
      expect_status(200)
      expect_json_types(
        id: :int,
        job_id: :int,
        note: :string_or_null,
        person_ids: :array_of_ints,
        scheduled_at: :date,
        length: :int_or_null,
        rating: :string_or_null,
        type: :string_or_null,
        created_at: :date,
        updated_at: :date
      )
    end
    it "should be able to show an offer" do
      get "#{@@base_url}/offers/#{@@offer_1_id}"
      expect_status(200)
      expect_json_types(
        id: :int,
        job_id: :int,
        note: :string_or_null,
        status: :string,
        created_at: :date,
        updated_at: :date
      )
    end
    it "should be able to show an status_change" do
      get "#{@@base_url}/status_changes/#{@@status_change_1_id}"
      expect_status(200)
      expect_json_types(
        id: :int,
        job_id: :int,
        note: :string_or_null,
        from: :string,
        to: :string,
        created_at: :date,
        updated_at: :date
      )
    end
    it "should be able to show a thanks_email" do
      get "#{@@base_url}/thanks_emails/#{@@thanks_email_1_id}"
      expect_status(200)
      expect_json_types(
        id: :int,
        job_id: :int,
        note: :string_or_null,
        person_ids: :array_of_ints,
        created_at: :date,
        updated_at: :date
      )
    end

  end
  describe "updates" do
    xit "should be able to update a company" do

    end
  end

  ## TODO 
  # describe "associations"
    # new job with person_ids
    # new person with job_id
    # each event with job
    # thanks_emails, followups, and interviews with person_ids
end
