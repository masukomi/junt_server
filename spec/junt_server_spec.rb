require 'airborne'

describe 'Junt' do
  before(:all) do
    ENV['JUNT_DB'] = "./test.db"
    # if ENV['JUNT_DB'].nil? 
    #   $stderr.puts("JUNT_DB env var must be specified.")
    #   exit 1
    # end
    # guarantee we're working with a fresh db
    `rm #{ENV['JUNT_DB']}`
    # start the junt server 
    @junt_server_pid = fork do
      exec "./junt"
    end
    puts "about to sleep..."
    # Process.detach(junt_server)
    # give the server a couple seconds to create the db
    # and start listening
    sleep(2) 


    @base_url = "http://localhost:8123/v1"
    @default_headers = {} # uses application/json content type by default
  end
  after(:all) do
    # Kill the server
    Process.kill("HUP", @junt_server_pid)
    # remove the test db
    `rm #{ENV['JUNT_DB']}`
  end
  describe "initial lists" do
    it "should initially have no companies" do
      get "#{@base_url}/companies" # should return empty array
      expect(json_body.size).to(eq(0))
    end
    it "should initially have no events" do
      get "#{@base_url}/events" # should return empty array
      expect(json_body.size).to(eq(0))
    end
    it "should initially have no followups" do
      get "#{@base_url}/followups" # should return empty array
      expect(json_body.size).to(eq(0))
    end
    it "should initially have no homeworks" do
      get "#{@base_url}/homeworks" # should return empty array
      expect(json_body.size).to(eq(0))
    end
    it "should initially have no interviews" do
      get "#{@base_url}/interviews" # should return empty array
      expect(json_body.size).to(eq(0))
    end
    it "should initially have no jobs" do
      get "#{@base_url}/jobs" # should return empty array
      expect(json_body.size).to(eq(0))
    end
    it "should initially have no offers" do
      get "#{@base_url}/offers" # should return empty array
      expect(json_body.size).to(eq(0))
    end
    it "should initially have no people" do
      get "#{@base_url}/people" # should return empty array
      expect(json_body.size).to(eq(0))
    end
    it "should initially have no status_changes" do
      get "#{@base_url}/status_changes" # should return empty array
      expect(json_body.size).to(eq(0))
    end
    it "should initially have no thanks_emails" do
      get "#{@base_url}/thanks_emails" # should return empty array
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
      post "#{@base_url}/companies", new_company_json, @default_headers
      expect_json_types(status: :string, id: :int)
      expect(json_body[:status]).to(eq("SUCCESS"))
      @company_1_id = json_body[:id]
    end
    it "should be able to create a job" do
      new_job_json = {
        company_id: @company_1_id,
        job_title: "chief tester",
        posting_url: "https://example.com/test_co/chief_tester_job",
        source: "we work remotely",
        referred_by: "nobody",
        salary_range: "unknown",
        application_method: "online",
        note: "a _markdown_ note",
        start_date: nil
      }
      post "#{@base_url}/jobs", new_job_json, @default_headers
      expect_json_types(status: :string, id: :int)
      expect(json_body[:status]).to(eq("SUCCESS"))
      @job_1_id = json_body[:id]
    end
    it "should be able to create a new person" do
      new_person_json = {
        name: "mary smith",
        email: 'msmith@example.com',
        phone: '+1 555-555-5555',
        note: 'another _markdown_ note',
      }
      post "#{@base_url}/people", new_person_json, @default_headers
      expect_json_types(status: :string, id: :int)
      expect(json_body[:status]).to(eq("SUCCESS"))
      @person_1_id = json_body[:id]
    end
    
  end

  ## TODO 
  # describe "associations"
    # new job with person_ids
    # new person with job_id

end
