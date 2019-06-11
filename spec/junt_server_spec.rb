require 'airborne'

class JuntTest
  attr_reader :server_pid, :base_url, :default_headers, :cached_data

  def initialize()
    # if ENV['JUNT_DB'].nil? 
      ENV['JUNT_DB'] = "./junt_test.db"
    # end
    @base_url="http://localhost:8123/v1"
    @default_headers = {}
    @cached_data = {}
  end
  def start()
    clean
    boot
  end
  def clean
    if File.exist? ENV['JUNT_DB']
      `rm #{ENV['JUNT_DB']}`
    end
  end
  def boot
    @server_pid = fork do
      exec "./junt_server"
    end
    puts "Server PID: #{@server_pid}"
    sleep(2)
  end
  def shutdown
    Process.kill("HUP", server_pid)
  end
  def self.begin()
    @@jt = JuntTest.new() unless defined?(@@jt)
    @@jt.start
    @jt
  end
  def self.end()
    @@jt.shutdown
  end
  def self.get()
    if ! defined?(@@jt)
      @@jt = JuntTest.new()
    end
    return @@jt
  end
end

describe 'Junt' do

  before(:all) do
    JuntTest.begin()
  end
  after(:all) do
    # Kill the server
    JuntTest.end()
  end
  describe "initial lists" do
    it "should initially have no companies" do
      get "#{JuntTest.get.base_url}/companies" # should return empty array
      expect(json_body.size).to(eq(0))
    end
    it "should initially have no events" do
      get "#{JuntTest.get.base_url}/events" # should return empty array
      expect(json_body.size).to(eq(0))
    end
    it "should initially have no followups" do
      get "#{JuntTest.get.base_url}/followups" # should return empty array
      expect(json_body.size).to(eq(0))
    end
    it "should initially have no homeworks" do
      get "#{JuntTest.get.base_url}/homeworks" # should return empty array
      expect(json_body.size).to(eq(0))
    end
    it "should initially have no interviews" do
      get "#{JuntTest.get.base_url}/interviews" # should return empty array
      expect(json_body.size).to(eq(0))
    end
    it "should initially have no jobs" do
      get "#{JuntTest.get.base_url}/jobs" # should return empty array
      expect(json_body.size).to(eq(0))
    end
    it "should initially have no offers" do
      get "#{JuntTest.get.base_url}/offers" # should return empty array
      expect(json_body.size).to(eq(0))
    end
    it "should initially have no people" do
      get "#{JuntTest.get.base_url}/people" # should return empty array
      expect(json_body.size).to(eq(0))
    end
    it "should initially have no status_changes" do
      get "#{JuntTest.get.base_url}/status_changes" # should return empty array
      expect(json_body.size).to(eq(0))
    end
    it "should initially have no thanks_emails" do
      get "#{JuntTest.get.base_url}/thanks_emails" # should return empty array
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
      post "#{JuntTest.get.base_url}/companies", new_company_json, JuntTest.get.default_headers
      expect_status(200)
      expect_json_types(status: :string, id: :int)
      expect(json_body[:status]).to(eq("SUCCESS"))
      JuntTest.get.cached_data[:company_1_id] = json_body[:id]
    end
    it "should be able to create a job" do
      new_job_json = {
        company_id: JuntTest.get.cached_data[:company_1_id],
        job_title: "chief tester",
        posting_url: "https://example.com/test_co/chief_tester_job",
        source: "we work remotely",
        referred_by: "nobody",
        salary_range: "unknown",
        application_method: "online",
        note: "a _markdown_ note",
        start_date: nil
      }
      post "#{JuntTest.get.base_url}/jobs", new_job_json, JuntTest.get.default_headers
      expect_status(200)
      expect_json_types(status: :string, id: :int)
      expect(json_body[:status]).to(eq("SUCCESS"))
      JuntTest.get.cached_data[:job_1_id] = json_body[:id]
    end
    it "should be able to create a new person" do
      new_person_json = {
        name: "mary smith",
        email: 'msmith@example.com',
        phone: '+1 555-555-5555',
        note: 'another _markdown_ note',
      }
      post "#{JuntTest.get.base_url}/people", new_person_json, JuntTest.get.default_headers
      expect_status(200)
      expect_json_types(status: :string, id: :int)
      expect(json_body[:status]).to(eq("SUCCESS"))
      JuntTest.get.cached_data[:person_1_id] = json_body[:id]
    end
    describe "of events" do
      it "should be able to create a new followup" do
        new_followup_json = {
          job_id: JuntTest.get.cached_data[:job_1_id],
          person_ids: [JuntTest.get.cached_data[:person_1_id]],
          note: "a _markdown_ note"
        }
        post "#{JuntTest.get.base_url}/followups", new_followup_json, JuntTest.get.default_headers
        expect_status(200)
        expect_json_types(status: :string, id: :int)
        expect(json_body[:status]).to(eq("SUCCESS"))
        JuntTest.get.cached_data[:followup_1_id] = json_body[:id]
      end
      it "should be able to create a new homework" do
        new_homework_json = {
          job_id: JuntTest.get.cached_data[:job_1_id],
          due_date: "2019-06-02T18:47:44+00:00",
          note: "a _markdown_ note"
        }
        post "#{JuntTest.get.base_url}/homeworks", new_homework_json, JuntTest.get.default_headers
        expect_status(200)
        expect_json_types(status: :string, id: :int)
        expect(json_body[:status]).to(eq("SUCCESS"))
        JuntTest.get.cached_data[:homework_1_id] = json_body[:id]
      end
      it "should be able to create a new interview" do
        new_interview_json = {
          job_id: JuntTest.get.cached_data[:job_1_id],
          person_ids: [JuntTest.get.cached_data[:person_1_id]],
          note: "a _markdown_ note",
          scheduled_at: "2019-06-02T18:47:44+00:00",
          rating: "😃",
          type: "technical"
        }
        post "#{JuntTest.get.base_url}/interviews", new_interview_json, JuntTest.get.default_headers
        expect_status(200)
        expect_json_types(status: :string, id: :int)
        expect(json_body[:status]).to(eq("SUCCESS"))
        JuntTest.get.cached_data[:interview_1_id] = json_body[:id]
      end
      it "should be able to create a new offer" do
        new_offer_json = {
          job_id: JuntTest.get.cached_data[:job_1_id],
          status: "received",
          note: "a _markdown_ note"
        }
        post "#{JuntTest.get.base_url}/offers", new_offer_json, JuntTest.get.default_headers
        expect_status(200)
        expect_json_types(status: :string, id: :int)
        expect(json_body[:status]).to(eq("SUCCESS"))
        JuntTest.get.cached_data[:offer_1_id] = json_body[:id]
      end
      it "should be able to create a new status_change" do
        new_status_change_json = {
          job_id: JuntTest.get.cached_data[:job_1_id],
          from: nil, 
          # initial status change for a job should be nil->something
          to: "applied",
          note: "a _markdown_ note"
        }
        post "#{JuntTest.get.base_url}/status_changes", new_status_change_json, JuntTest.get.default_headers
        expect_status(200)
        expect_json_types(status: :string, id: :int)
        expect(json_body[:status]).to(eq("SUCCESS"))
        JuntTest.get.cached_data[:status_change_1_id] = json_body[:id]
      end
      it "should be able to create a new thanks_email" do
        new_thanks_email_json = {
          job_id: JuntTest.get.cached_data[:job_1_id],
          person_ids: [JuntTest.get.cached_data[:person_1_id]],
          note: "a _markdown_ note" 
        }
        post "#{JuntTest.get.base_url}/thanks_emails", new_thanks_email_json, JuntTest.get.default_headers
        expect_status(200)
        expect_json_types(status: :string, id: :int)
        expect(json_body[:status]).to(eq("SUCCESS"))
        JuntTest.get.cached_data[:thanks_email_1_id] = json_body[:id]
      end
    end # end events creation
    
  end # end creation
  describe "showing objects" do
    it "should be able to show a company" do
      get "#{JuntTest.get.base_url}/companies/#{JuntTest.get.cached_data[:company_1_id]}"
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
      get "#{JuntTest.get.base_url}/jobs/#{JuntTest.get.cached_data[:job_1_id]}"
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
      get "#{JuntTest.get.base_url}/people/#{JuntTest.get.cached_data[:person_1_id]}"
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
      get "#{JuntTest.get.base_url}/followups/#{JuntTest.get.cached_data[:followup_1_id]}"
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
      get "#{JuntTest.get.base_url}/homeworks/#{JuntTest.get.cached_data[:homework_1_id]}"
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
      get "#{JuntTest.get.base_url}/interviews/#{JuntTest.get.cached_data[:interview_1_id]}"
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
      get "#{JuntTest.get.base_url}/offers/#{JuntTest.get.cached_data[:offer_1_id]}"
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
      get "#{JuntTest.get.base_url}/status_changes/#{JuntTest.get.cached_data[:status_change_1_id]}"
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
      get "#{JuntTest.get.base_url}/thanks_emails/#{JuntTest.get.cached_data[:thanks_email_1_id]}"
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
