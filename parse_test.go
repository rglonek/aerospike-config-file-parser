package aeroconf

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestParseBasic(t *testing.T) {
	items := `# Aerospike database configuration file for use with systemd.

	service {
		paxos-single-replica-limit 1 # Number of nodes where the replica count is automatically reduced to 1.
		proto-fd-max 15000
	}
	
	logging {
		file /var/log/aerospike.log {
			context any info
		}
		file /var/log/aerospike-ldap.log {
			context security detail
		}
	}
	
	security {
		enable-security true
		enable-ldap true
		ldap {
			query-base-dn dc=aerospike,dc=com
			server ldap://LDAPIP:389
			disable-tls true
			user-dn-pattern uid=${un},ou=People,dc=aerospike,dc=com
			role-query-search-ou false
			role-query-pattern (&(objectClass=posixGroup)(memberUid=${un}))
			polling-period 10
		}
		log {
			report-violation true
		}
	}
	
	network {
		tls tls1 {
		cert-file /etc/aerospike/ssl/tls1/cert.pem
		key-file /etc/aerospike/ssl/tls1/key.pem
		ca-file /etc/aerospike/ssl/tls1/cacert.pem
		}
		service {
			#address any
			#port 3000
	
			tls-port 4333
			tls-address any
			tls-authenticate-client false
			# could be any | user-defined
			tls-name tls1
		}
	
		heartbeat {
			mode multicast
			multicast-group 239.1.99.222
			port 9918
	
			# To use unicast-mesh heartbeats, remove the 3 lines above, and see
			# aerospike_mesh.conf for alternative.
	
			interval 150
			timeout 10
		}
	
		fabric {
			#port 3001
			tls-port 3011
			tls-name tls1
		}
	
		info {
			port 3003
		}
	}
	
	namespace test {
		replication-factor 2
		memory-size 4G
		default-ttl 30d # 30 days, use 0 to never expire/evict.
	
		storage-engine memory
	}
	
	namespace bar {
		replication-factor 2
		memory-size 4G
		default-ttl 30d # 30 days, use 0 to never expire/evict.
	
		storage-engine memory
	
		# To use file storage backing, comment out the line above and use the
		# following lines instead.
	#	storage-engine device {
	#		file /opt/aerospike/data/bar.dat
	#		filesize 16G
	#		data-in-memory true # Store data in memory in addition to file.
	#	}
	}
	`
	r := strings.NewReader(items)
	s, err := Parse(r)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	fmt.Println(s.Type("service"))
	fmt.Println(s.Stanza("service").Type("proto-fd-max"))
	out, _ := s.Stanza("service").GetValues("proto-fd-max")
	for _, i := range out {
		fmt.Println(*i)
	}
	if s.Stanza("service") == nil {
		s.NewStanza("service")
	}
	s.Stanza("service").SetValue("proto-fd-max", "30000")

	if s.Stanza("network") == nil {
		s.NewStanza("network")
	}
	if s.Stanza("network").Stanza("heartbeat") == nil {
		s.Stanza("network").NewStanza("heartbeat")
	}
	s.Stanza("network").Stanza("heartbeat").Delete("multicast-group")
	s.Stanza("network").Stanza("heartbeat").SetValue("mode", "mesh")
	s.Stanza("network").Stanza("heartbeat").SetValues("mesh-seed-address-port", SliceToValues([]string{"172.17.0.2 3000", "172.17.0.3 3000"}))

	s.Stanza("network").Delete("info")
	s.Stanza("network").NewStanza("info")
	s.Stanza("network").Stanza("info").SetValue("port", "3003")
	err = s.Write(os.Stdout, "", "    ", true)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
}
